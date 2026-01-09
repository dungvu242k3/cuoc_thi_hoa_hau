import { createFileRoute } from "@tanstack/react-router";
import { ContestantLayout } from "../pages/ContestantLayout";

export const Route = createFileRoute("/contestant")({
    component: ContestantLayout,
});
