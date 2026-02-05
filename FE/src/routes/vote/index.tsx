import { createFileRoute } from '@tanstack/react-router';

function VoteDashboard() {
    return (
        <div className="w-full">
            {/* Hero Section with Poster */}
            <div className="relative w-full overflow-hidden">
                <div className="max-w-6xl mx-auto px-4 py-8">
                    <div className="relative rounded-2xl overflow-hidden shadow-2xl">
                        <img
                            src="/poster.jpeg"
                            alt="Concert Việt Nam Bùng Sáng 2026"
                            className="w-full h-auto object-cover"
                        />
                        <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent"></div>
                    </div>
                </div>
            </div>
            {/* Main Content - Empty as requested */}
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16 flex-1">
            </main>
        </div>
    );
}

export const Route = createFileRoute('/vote/')({
    component: VoteDashboard,
});
