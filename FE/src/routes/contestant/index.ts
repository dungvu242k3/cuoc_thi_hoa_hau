import { createFileRoute } from "@tanstack/react-router";
import { ContestantDispatcher } from "../../pages/ContestantDispatcher";

export const Route = createFileRoute("/contestant/")({
  component: ContestantDispatcher,
});
