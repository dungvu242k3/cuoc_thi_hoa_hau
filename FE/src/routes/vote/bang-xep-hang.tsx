import { useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { gql } from 'graphql-request';
import { useMemo } from 'react';
import { apiRequest } from '../../lib/api';

const GET_PUBLIC_CONTESTANTS = gql`
  query PublicContestants($page: Int!, $limit: Int!) {
    publicContestants(page: $page, limit: $limit) {
      id
      sbd
      voteCount
      personalInfo {
        fullName
        dob
        nationality
        address
        job
      }
      portfolio {
        avatarUrl
      }
    }
  }
`;

function Leaderboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['publicContestants'],
    queryFn: async () => apiRequest(GET_PUBLIC_CONTESTANTS, { page: 1, limit: 100 }),
  });

  const sortedContestants = useMemo(() => {
    if (!data?.publicContestants) return [];
    return [...data.publicContestants]
      .sort((a: any, b: any) => (b.voteCount || 0) - (a.voteCount || 0))
      .map((c: any, index: number) => ({ ...c, rank: index + 1 }));
  }, [data?.publicContestants]);

  const top3 = sortedContestants.slice(0, 3);
  const others = sortedContestants.slice(3);

  if (isLoading) {
    return (
      <div className="min-h-[50vh] flex items-center justify-center">
        <div className="text-white text-2xl font-bold animate-pulse">Đang tải...</div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Header */}
      <div className="text-center mb-10 md:mb-16">
        <h1 className="text-3xl md:text-5xl font-bold text-white mb-2 md:mb-3 tracking-wide drop-shadow-lg">
          BẢNG XẾP HẠNG
        </h1>
        <p className="text-lg md:text-xl text-white/90 font-light">
          CẬP NHẬT THEO THỜI GIAN THỰC
        </p>
      </div>

      {/* Top 3 Podium */}
      {top3.length > 0 && (
        <div className="flex flex-row items-end justify-center mb-12 md:mb-20 gap-3 md:gap-8 px-2 md:px-0">
          {/* Rank 2 - Left */}
          {top3[1] && (
            <div className="relative flex flex-col items-center order-1 w-1/3 md:w-auto">
              <div className="mb-2 md:mb-4 relative">
                <div className="w-16 h-16 md:w-32 md:h-32 rounded-full overflow-hidden border-2 md:border-4 border-gray-300 shadow-[0_0_15px_rgba(255,255,255,0.2)]">
                  <img src={top3[1].portfolio?.avatarUrl || "https://placehold.co/150x200"} alt={top3[1].personalInfo?.fullName} className="w-full h-full object-cover" />
                </div>
                <div className="absolute -bottom-2 md:-bottom-3 left-1/2 -translate-x-1/2 bg-gray-300 text-black font-bold w-5 h-5 md:w-8 md:h-8 rounded-full flex items-center justify-center border md:border-2 border-white text-xs md:text-base">2</div>
              </div>
              <div className="text-center w-full">
                <div className="text-white font-bold text-xs md:text-lg truncate px-1">{top3[1].personalInfo?.fullName}</div>
                <div className="text-white/80 text-[10px] md:text-sm hidden md:block">SBD: {top3[1].sbd}</div>
                <div className="text-yellow-400 font-bold text-xs md:text-xl mt-0.5">{top3[1].voteCount} <span className="md:hidden">vote</span></div>
              </div>
            </div>
          )}

          {/* Rank 1 - Center */}
          {top3[0] && (
            <div className="relative flex flex-col items-center order-2 w-1/3 md:w-auto -mb-4 md:mb-0 z-10 scale-110 md:scale-100 origin-bottom">
              <div className="mb-4 md:mb-6 relative">
                <div className="absolute -top-8 md:-top-12 left-1/2 -translate-x-1/2 text-yellow-400">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 md:h-12 md:w-12" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                </div>
                <div className="w-24 h-24 md:w-40 md:h-40 rounded-full overflow-hidden border-2 md:border-4 border-yellow-400 shadow-[0_0_25px_rgba(250,204,21,0.4)]">
                  <img src={top3[0].portfolio?.avatarUrl || "https://placehold.co/150x200"} alt={top3[0].personalInfo?.fullName} className="w-full h-full object-cover" />
                </div>
                <div className="absolute -bottom-3 md:-bottom-4 left-1/2 -translate-x-1/2 bg-yellow-400 text-black font-bold w-6 h-6 md:w-10 md:h-10 rounded-full flex items-center justify-center border md:border-2 border-white text-sm md:text-xl">1</div>
              </div>
              <div className="text-center w-full">
                <div className="text-xl md:text-2xl text-white font-bold leading-tight truncate px-1">{top3[0].personalInfo?.fullName}</div>
                <div className="text-white/80 text-[10px] md:text-base hidden md:block">SBD: {top3[0].sbd}</div>
                <div className="text-yellow-400 font-bold text-sm md:text-3xl mt-1">{top3[0].voteCount} <span className="md:hidden">vote</span></div>
              </div>
            </div>
          )}

          {/* Rank 3 - Right */}
          {top3[2] && (
            <div className="relative flex flex-col items-center order-3 w-1/3 md:w-auto">
              <div className="mb-2 md:mb-4 relative">
                <div className="w-16 h-16 md:w-32 md:h-32 rounded-full overflow-hidden border-2 md:border-4 border-amber-700 shadow-[0_0_15px_rgba(180,83,9,0.2)]">
                  <img src={top3[2].portfolio?.avatarUrl || "https://placehold.co/150x200"} alt={top3[2].personalInfo?.fullName} className="w-full h-full object-cover" />
                </div>
                <div className="absolute -bottom-2 md:-bottom-3 left-1/2 -translate-x-1/2 bg-amber-700 text-white font-bold w-5 h-5 md:w-8 md:h-8 rounded-full flex items-center justify-center border md:border-2 border-white text-xs md:text-base">3</div>
              </div>
              <div className="text-center w-full">
                <div className="text-white font-bold text-xs md:text-lg truncate px-1">{top3[2].personalInfo?.fullName}</div>
                <div className="text-white/80 text-[10px] md:text-sm hidden md:block">SBD: {top3[2].sbd}</div>
                <div className="text-yellow-400 font-bold text-xs md:text-xl mt-0.5">{top3[2].voteCount} <span className="md:hidden">vote</span></div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* The Rest List - Desktop Table & Mobile Cards */}
      <div className="bg-white/10 backdrop-blur-md rounded-3xl overflow-hidden border border-white/10 shadow-xl">
        {/* Desktop View: Table */}
        <div className="hidden md:block overflow-x-auto">
          <table className="w-full text-white min-w-[600px]">
            <thead>
              <tr className="bg-black/20 text-left">
                <th className="py-3 md:py-4 px-3 md:px-6 font-semibold w-16 md:w-24 text-center">Hạng</th>
                <th className="py-3 md:py-4 px-3 md:px-6 font-semibold">Thí sinh</th>
                <th className="py-3 md:py-4 px-3 md:px-6 font-semibold w-24 md:w-32 text-center">SBD</th>
                <th className="py-3 md:py-4 px-3 md:px-6 font-semibold w-32 md:w-40 text-right">Vote</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {others.map((c: any) => (
                <tr key={c.id} className="hover:bg-white/5 transition-colors">
                  <td className="py-3 md:py-4 px-3 md:px-6 text-center font-bold text-white/50 text-lg md:text-xl">{c.rank}</td>
                  <td className="py-3 md:py-4 px-3 md:px-6">
                    <div className="flex items-center space-x-2 md:space-x-4">
                      <img src={c.portfolio?.avatarUrl || "https://placehold.co/100"} alt="" className="w-10 h-10 md:w-12 md:h-12 rounded-full object-cover border border-white/20" />
                      <div>
                        <div className="font-semibold text-sm md:text-base">{c.personalInfo?.fullName}</div>
                        <div className="text-xs text-white/60">{c.personalInfo?.address}</div>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 md:py-4 px-3 md:px-6 text-center font-mono text-white/80 text-sm md:text-base">{c.sbd}</td>
                  <td className="py-3 md:py-4 px-3 md:px-6 text-right font-bold text-yellow-400 text-sm md:text-base">{c.voteCount}</td>
                </tr>
              ))}
              {sortedContestants.length === 0 && (
                <tr>
                  <td colSpan={4} className="py-8 text-center text-white/50">Chưa có dữ liệu</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Mobile View: Cards List */}
        <div className="md:hidden divide-y divide-white/10">
          {others.map((c: any) => (
            <div key={c.id} className="p-4 flex items-center gap-4 hover:bg-white/5 active:bg-white/10 transition">
              <div className="text-xl font-bold text-white/40 w-8 text-center">{c.rank}</div>
              <img
                src={c.portfolio?.avatarUrl || "https://placehold.co/100"}
                alt=""
                className="w-12 h-12 rounded-full object-cover border border-white/20 shadow-sm"
              />
              <div className="flex-1 min-w-0">
                <div className="font-bold text-white text-base truncate">{c.personalInfo?.fullName}</div>
                <div className="text-xs text-white/60 flex items-center gap-2">
                  <span className="bg-white/10 px-1.5 py-0.5 rounded text-[10px] font-mono">SBD: {c.sbd}</span>
                  <span className="truncate">{c.personalInfo?.address}</span>
                </div>
              </div>
              <div className="text-right">
                <div className="font-bold text-yellow-400 text-lg">{c.voteCount}</div>
                <div className="text-[10px] text-white/50 uppercase tracking-wider">Vote</div>
              </div>
            </div>
          ))}
          {sortedContestants.length === 0 && (
            <div className="py-10 text-center text-white/50 italic">
              Chưa có dữ liệu xếp hạng
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/vote/bang-xep-hang')({
  component: Leaderboard,
});
