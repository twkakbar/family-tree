import { ExploreExperience } from "@/components/tree/explore-experience";

export const metadata = {
  title: "Kesultanan Aceh",
  description: "Explore the Aceh royal family tree in 3D.",
};

/** Public entry: the immersive tree is the first page. */
export default function HomePage() {
  return <ExploreExperience />;
}
