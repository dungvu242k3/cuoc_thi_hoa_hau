
import { createFileRoute } from '@tanstack/react-router';
import { VoteAuthComponent } from './-AuthComponent.tsx';

export const Route = createFileRoute('/vote/auth')({
    component: VoteAuthComponent,
});
