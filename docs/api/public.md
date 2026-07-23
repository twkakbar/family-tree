# Public API

Unauthenticated reads of the singleton Aceh royal family tree.
Base: `/api/v1/public`.

## `GET /public/meta`

Family branding metadata.

**Response `200`**

```json
{ "id": "…", "name": "Kesultanan Aceh", "description": "…" }
```

## `GET /public/tree`

Full graph for the 3D explorer.

**Response `200`**

```json
{
  "persons": [ /* Person objects — see persons.md */ ],
  "relationships": [ /* Relationship objects */ ]
}
```

## `GET /public/persons/:id`

Person detail including incident relationships. Returns `404` if the person
is not part of the royal family.

## `GET /public/persons/:id/ancestors?depth=5`

## `GET /public/persons/:id/descendants?depth=5`

Same shapes as the authenticated traversal endpoints in [`persons.md`](persons.md).
