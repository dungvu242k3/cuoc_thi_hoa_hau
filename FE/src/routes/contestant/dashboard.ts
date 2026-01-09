import { createFileRoute } from "@tanstack/react-router";
import { ContestantDashboard } from "../../pages/ContestantDashboard";

export const Route = createFileRoute("/contestant/dashboard")({
    component: ContestantDashboard,
});
