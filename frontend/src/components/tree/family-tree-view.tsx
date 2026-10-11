"use client";

import { useCallback, useMemo, useRef, useState } from "react";

import { Card, CardDescription, CardTitle } from "@/components/ui/card";
import { Spinner } from "@/components/ui/spinner";
import {
  RelationshipLinkPicker,
  choiceToRelationshipInput,
  type RelationshipLinkChoice,
} from "@/components/tree/relationship-link-picker";
import { PersonTreeNode } from "@/components/tree/person-tree-node";
import { TreeToolbar } from "@/components/tree/tree-toolbar";
import { useCreateRelationship, useFamilyTree } from "@/hooks/use-persons";
import { ApiError } from "@/types/api";
import type { Person, TreeNodePosition } from "@/types/person";
import {
  edgePath,
  findNode,
  layoutFamilyTree,
  nodeCenter,
  treeMetrics,
} from "@/utils/tree-layout";
import { cn } from "@/utils/cn";

const LINK_DRAG_THRESHOLD_PX = 8;
const MIN_SCALE = 0.35;
const MAX_SCALE = 2;

interface FamilyTreeViewProps {
  familyId: string;
  canEdit?: boolean;
  selectedPersonId?: string | null;
  onSelectPerson?: (person: Person) => void;
}

interface LinkDragState {
  sourcePerson: Person;
  sourceNode: TreeNodePosition;
  cursor: { x: number; y: number };
  hoverTargetId: string | null;
}

interface PendingLink {
  source: Person;
  target: Person;
}

/** Pan/zoom tree with drag-to-link relationship creation for editors. */
export function FamilyTreeView({
  familyId,
  canEdit = false,
  selectedPersonId,
  onSelectPerson,
}: FamilyTreeViewProps) {
  const tree = useFamilyTree(familyId);
  const createRelationship = useCreateRelationship(familyId);
  const containerRef = useRef<HTMLDivElement>(null);
  const [scale, setScale] = useState(1);
  const [offset, setOffset] = useState({ x: 0, y: 0 });
  const [isPanning, setIsPanning] = useState(false);
  const panRef = useRef({ x: 0, y: 0 });
  const panningRef = useRef(false);
  const [linkDrag, setLinkDrag] = useState<LinkDragState | null>(null);
  const [pendingLink, setPendingLink] = useState<PendingLink | null>(null);

  const layout = useMemo(
    () => (tree.data ? layoutFamilyTree(tree.data) : null),
    [tree.data],
  );

  const clampScale = (value: number) => Math.min(MAX_SCALE, Math.max(MIN_SCALE, value));

  const resetView = useCallback(() => {
    setScale(1);
    setOffset({ x: 0, y: 0 });
  }, []);

  const clientToLayout = useCallback(
    (clientX: number, clientY: number) => {
      const rect = containerRef.current?.getBoundingClientRect();
      if (!rect) return { x: 0, y: 0 };
      return {
        x: (clientX - rect.left - offset.x) / scale,
        y: (clientY - rect.top - offset.y) / scale,
      };
    },
    [offset.x, offset.y, scale],
  );

  const hitTestNode = useCallback(
    (clientX: number, clientY: number): TreeNodePosition | null => {
      if (!layout) return null;
      const point = clientToLayout(clientX, clientY);
      for (const node of layout.nodes) {
        if (
          point.x >= node.x &&
          point.x <= node.x + treeMetrics.NODE_WIDTH &&
          point.y >= node.y &&
          point.y <= node.y + treeMetrics.NODE_HEIGHT
        ) {
          return node;
        }
      }
      return null;
    },
    [clientToLayout, layout],
  );

  const startLinkDrag = useCallback(
    (event: React.PointerEvent, person: Person, node: TreeNodePosition) => {
      if (!canEdit) return;
      event.stopPropagation();
      event.preventDefault();

      const startX = event.clientX;
      const startY = event.clientY;
      let linking = false;

      const onPointerMove = (ev: PointerEvent) => {
        const dx = ev.clientX - startX;
        const dy = ev.clientY - startY;
        if (!linking && Math.hypot(dx, dy) < LINK_DRAG_THRESHOLD_PX) return;

        linking = true;
        const cursor = clientToLayout(ev.clientX, ev.clientY);
        const hover = hitTestNode(ev.clientX, ev.clientY);
        setLinkDrag({
          sourcePerson: person,
          sourceNode: node,
          cursor,
          hoverTargetId:
            hover && hover.person.id !== person.id ? hover.person.id : null,
        });
      };

      const onPointerUp = (ev: PointerEvent) => {
        window.removeEventListener("pointermove", onPointerMove);
        window.removeEventListener("pointerup", onPointerUp);

        if (linking) {
          const targetNode = hitTestNode(ev.clientX, ev.clientY);
          if (targetNode && targetNode.person.id !== person.id) {
            setPendingLink({ source: person, target: targetNode.person });
          }
          setLinkDrag(null);
        } else {
          onSelectPerson?.(person);
        }
      };

      window.addEventListener("pointermove", onPointerMove);
      window.addEventListener("pointerup", onPointerUp);
    },
    [canEdit, clientToLayout, hitTestNode, onSelectPerson],
  );

  const handleConfirmLink = (choice: RelationshipLinkChoice) => {
    if (!pendingLink) return;
    const input = choiceToRelationshipInput(choice, pendingLink.source, pendingLink.target);
    createRelationship.mutate(input, {
      onSuccess: () => {
        createRelationship.reset();
        setPendingLink(null);
      },
    });
  };

  const pickerError =
    createRelationship.error instanceof ApiError && !createRelationship.isPending
      ? createRelationship.error.message
      : null;

  if (tree.isPending) {
    return (
      <Card className="grid min-h-[min(24rem,50vh)] place-items-center">
        <Spinner className="size-8 text-emerald-600" />
      </Card>
    );
  }

  if (!layout || layout.nodes.length === 0) {
    return (
      <Card className="p-4 sm:p-6">
        <CardTitle>Family tree</CardTitle>
        <CardDescription className="mt-1">
          Add members below, then drag one card onto another to connect them.
        </CardDescription>
      </Card>
    );
  }

  return (
    <>
      <Card className="overflow-hidden p-0">
        <TreeToolbar
          canEdit={canEdit}
          scale={scale}
          onZoomIn={() => setScale((value) => clampScale(value + 0.15))}
          onZoomOut={() => setScale((value) => clampScale(value - 0.15))}
          onResetView={resetView}
        />

        <div
          ref={containerRef}
          className={cn(
            "relative min-h-[min(22rem,55vh)] h-[min(28rem,62vh)] w-full overflow-hidden bg-[radial-gradient(circle_at_1px_1px,#e2e8f0_1px,transparent_0)] [background-size:20px_20px] sm:min-h-[26rem] sm:h-[32rem]",
            "touch-none",
            linkDrag ? "cursor-crosshair" : isPanning ? "cursor-grabbing" : "cursor-grab",
          )}
          onWheel={(event) => {
            event.preventDefault();
            const delta = event.deltaY > 0 ? -0.08 : 0.08;
            setScale((value) => clampScale(value + delta));
          }}
          onPointerDown={(event) => {
            if (event.pointerType === "mouse" && event.button !== 0) return;
            if ((event.target as Element).closest("[data-person-node]")) return;
            panRef.current = {
              x: event.clientX - offset.x,
              y: event.clientY - offset.y,
            };
            panningRef.current = true;
            setIsPanning(true);
            event.currentTarget.setPointerCapture(event.pointerId);
          }}
          onPointerMove={(event) => {
            if (!panningRef.current || linkDrag) return;
            setOffset({
              x: event.clientX - panRef.current.x,
              y: event.clientY - panRef.current.y,
            });
          }}
          onPointerUp={() => {
            panningRef.current = false;
            setIsPanning(false);
          }}
          onPointerCancel={() => {
            panningRef.current = false;
            setIsPanning(false);
          }}
        >
          <div
            className="absolute left-0 top-0 origin-top-left will-change-transform"
            style={{
              transform: `translate(${offset.x}px, ${offset.y}px) scale(${scale})`,
            }}
          >
            <svg width={layout.width} height={layout.height} className="overflow-visible">
              {layout.edges.map((edge) => {
                const from = findNode(layout.nodes, edge.from_person_id);
                const to = findNode(layout.nodes, edge.to_person_id);
                if (!from || !to) return null;
                return (
                  <path
                    key={edge.id}
                    d={edgePath(from, to, edge.type)}
                    fill="none"
                    stroke={edge.type === "spouse" ? "#94a3b8" : "#10b981"}
                    strokeWidth={edge.type === "spouse" ? 2 : 2.5}
                    strokeDasharray={edge.type === "spouse" ? "6 4" : undefined}
                  />
                );
              })}

              {linkDrag && (
                <>
                  <line
                    x1={nodeCenter(linkDrag.sourceNode).x}
                    y1={nodeCenter(linkDrag.sourceNode).y}
                    x2={linkDrag.cursor.x}
                    y2={linkDrag.cursor.y}
                    stroke="#059669"
                    strokeWidth={2.5}
                    strokeDasharray="6 4"
                  />
                  {linkDrag.hoverTargetId && (() => {
                    const target = findNode(layout.nodes, linkDrag.hoverTargetId);
                    if (!target) return null;
                    const center = nodeCenter(target);
                    return (
                      <circle
                        cx={center.x}
                        cy={center.y}
                        r={14}
                        fill="#10b981"
                        fillOpacity={0.2}
                        stroke="#059669"
                        strokeWidth={2}
                      />
                    );
                  })()}
                </>
              )}
            </svg>

            {layout.nodes.map((node) => (
              <PersonTreeNode
                key={node.person.id}
                node={node}
                canEdit={canEdit}
                isSource={linkDrag?.sourcePerson.id === node.person.id}
                isHoverTarget={linkDrag?.hoverTargetId === node.person.id}
                isSelected={selectedPersonId === node.person.id}
                isLinking={Boolean(linkDrag)}
                width={treeMetrics.NODE_WIDTH}
                height={treeMetrics.NODE_HEIGHT}
                onPointerDown={(event) => startLinkDrag(event, node.person, node)}
                onSelect={() => onSelectPerson?.(node.person)}
              />
            ))}
          </div>
        </div>
      </Card>

      {pendingLink && (
        <RelationshipLinkPicker
          source={pendingLink.source}
          target={pendingLink.target}
          relationships={tree.data?.relationships ?? []}
          loading={createRelationship.isPending}
          error={pickerError}
          onChoose={handleConfirmLink}
          onCancel={() => {
            createRelationship.reset();
            setPendingLink(null);
          }}
        />
      )}
    </>
  );
}
