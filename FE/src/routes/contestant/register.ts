import { createFileRoute } from "@tanstack/react-router";
import { ContestantRegistration } from "../../pages/ContestantRegistration";

export const Route = createFileRoute("/contestant/register")({
    component: ContestantRegistration,
});
