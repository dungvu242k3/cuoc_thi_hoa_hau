import { useMutation, useQuery } from '@tanstack/react-query';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { gql } from 'graphql-request';
import { useMemo, useState } from 'react';
import { apiRequest } from '../../lib/api';
import { useAuthStore } from '../../store/useAuthStore';

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
      physicalInfo {
        height
        weight
        measurements
      }
      portfolio {
        avatarUrl
      }
    }
  }
`;

function ContestantList() {
  const navigate = useNavigate();
  const { token, role, isAuthenticated } = useAuthStore();
  const [searchTerm, setSearchTerm] = useState('');

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['publicContestants'],
    queryFn: async () => apiRequest(GET_PUBLIC_CONTESTANTS, { page: 1, limit: 100 }),
  });

  const sortedContestants = useMemo(() => {
    if (!data?.publicContestants) return [];

    let filtered = [...data.publicContestants];
    if (searchTerm) {
      const lowerTerm = searchTerm.toLowerCase();
      filtered = filtered.filter((c: any) =>
        (c.personalInfo?.fullName?.toLowerCase().includes(lowerTerm)) ||
        (c.sbd?.toLowerCase().includes(lowerTerm))
      );
    }

    return filtered
      .sort((a: any, b: any) => (b.voteCount || 0) - (a.voteCount || 0))
      .map((c: any, index: number) => ({ ...c, rank: index + 1 }));
  }, [data?.publicContestants, searchTerm]);

  const voteMutation = useMutation({
    mutationFn: async (contestantId: string) => {
      if (!isAuthenticated || !token) {
        navigate({ to: '/vote/auth' as '/vote/auth' });
        throw new Error("Unauthorized");
      }

      const voteResult = await apiRequest(
        `mutation VoteForContestant($id: ID!) { voteForContestant(id: $id) }`,
        { id: contestantId },
        token
      );
      return voteResult.voteForContestant;
    },
    onSuccess: () => {
      refetch();
      alert("Bình chọn thành công!");
    },
    onError: (err: any) => {
      if (err.message !== "Unauthorized") {
        alert(err.message || "Bình chọn thất bại");
      }
    }
  });

  const handleVote = (id: string) => {
    if (!isAuthenticated) {
      navigate({ to: '/vote/auth' as '/vote/auth' });
      return;
    }
    if (role !== 'audience' && role !== 'admin') {
      alert("Chỉ tài khoản khán giả mới có thể bình chọn!");
      return;
    }
    voteMutation.mutate(id);
  };

  if (isLoading) {
    return (
      <div className="min-h-[50vh] flex items-center justify-center">
        <div className="text-white text-2xl font-bold animate-pulse">Đang tải...</div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Page Header */}
      <div className="text-center mb-8 md:mb-12">
        <h1 className="text-3xl md:text-5xl font-bold text-white mb-2 md:mb-3 tracking-wide">
          DANH SÁCH THÍ SINH
        </h1>
        <p className="text-lg md:text-xl text-white/90">
          CUỘC THI HOA HẬU VIỆT NAM 2026
        </p>
      </div>

      {/* Search Bar */}
      <div className="max-w-2xl mx-auto mb-8 md:mb-12 relative">
        <input
          type="text"
          placeholder="Tìm kiếm theo Tên hoặc SBD..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full py-3 md:py-4 px-5 md:px-6 pr-12 rounded-full bg-white/10 backdrop-blur-md border border-white/20 text-white placeholder-white/60 focus:outline-none focus:ring-2 focus:ring-white/50 text-base md:text-lg shadow-lg"
        />
        <div className="absolute right-4 top-1/2 -translate-y-1/2 text-white/60">
          <svg className="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>

      {/* Contestants Grid - Responsive Gap */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-8 mb-12">
        {sortedContestants.map((c: any) => (
          <div key={c.id} className="w-full">
            {/* DESKTOP VIEW */}
            <div className="hidden md:block h-full bg-gradient-to-b from-blue-300 to-blue-600 rounded-3xl overflow-hidden shadow-xl hover:shadow-2xl hover:scale-[1.02] transition-all duration-300 p-1">
              <div className="bg-white/10 backdrop-blur-sm rounded-[22px] h-full flex flex-col p-3">
                {/* Image Container */}
                <div className="relative aspect-[3/4] w-full overflow-hidden rounded-2xl mb-3 group">
                  <img
                    src={c.portfolio?.avatarUrl || "https://placehold.co/300x400?text=No+Image"}
                    alt={c.personalInfo?.fullName || "Thí sinh"}
                    loading="lazy"
                    className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                  />
                </div>

                {/* SBD & Vote Count Bar - Moved Below Image */}
                <div className="mb-3 h-12 bg-white/20 backdrop-blur-md rounded-xl flex items-center justify-between px-4 border border-white/30">
                  <div className="text-white font-bold text-sm tracking-wider">
                    SBD: {c.sbd}
                  </div>
                  <div className="text-white font-bold text-lg">
                    {c.voteCount || 0}
                  </div>
                </div>

                {/* Info Section */}
                <div className="px-2 pb-2 flex-1 flex flex-col">
                  <div className="text-white/90 text-sm font-semibold mb-1 uppercase tracking-wide">
                    {c.personalInfo?.address || c.personalInfo?.nationality || "N/A"}
                  </div>

                  <h3 className="text-white text-xl font-extrabold uppercase leading-tight mb-4 tracking-wide shadow-black drop-shadow-sm line-clamp-2 min-h-[3rem]">
                    {c.personalInfo?.fullName}
                  </h3>

                  <div className="mt-auto flex items-end justify-between">
                    {/* Vote Button */}
                    <button
                      onClick={() => handleVote(c.id)}
                      disabled={voteMutation.isPending}
                      className="bg-white text-blue-600 font-extrabold px-8 py-3 rounded-full hover:bg-blue-50 transition-colors shadow-lg disabled:opacity-70 disabled:cursor-not-allowed"
                    >
                      {voteMutation.isPending ? '...' : 'Bình chọn'}
                    </button>

                    {/* Rank Badge (Laurel Wreath) */}
                    <div className="relative w-16 h-16 flex-shrink-0 text-white/90">
                      <svg viewBox="0 0 100 100" fill="currentColor" className="w-full h-full drop-shadow-md">
                        <path d="M50 10 C20 10 5 35 5 60 C5 75 15 85 25 90 L25 85 C18 80 10 70 10 60 C10 40 25 15 50 15 C75 15 90 40 90 60 C90 70 82 80 75 85 L75 90 C85 85 95 75 95 60 C95 35 80 10 50 10 Z" />
                        <path d="M50 85 L45 95 L50 90 L55 95 Z" />
                        <path d="M15 60 Q15 45 30 35 T50 30 T70 35 T85 60" fill="none" stroke="currentColor" strokeWidth="2" strokeDasharray="2 2" />
                      </svg>
                      <div className="absolute inset-0 flex items-center justify-center pt-2">
                        <span className="font-bold text-xl">{c.rank}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* MOBILE VIEW - COMPACT HORIZONTAL CARD */}
            <div className="md:hidden bg-gradient-to-r from-white/10 to-white/5 backdrop-blur-md border border-white/20 rounded-2xl p-3 flex gap-4 shadow-lg hover:bg-white/10 transition-colors">
              {/* Left: Avatar (Square/Rounded) */}
              <div className="relative w-24 aspect-[3/4] flex-shrink-0">
                <img
                  src={c.portfolio?.avatarUrl || "https://placehold.co/300x400?text=No+Image"}
                  alt={c.personalInfo?.fullName}
                  className="w-full h-full object-cover rounded-xl shadow-md border border-white/10"
                  loading="lazy"
                />
                {/* Rank Badge Overlay */}
                <div className="absolute -top-1 -left-1 bg-yellow-400 text-blue-900 text-xs font-bold w-6 h-6 rounded-full flex items-center justify-center shadow-sm border border-white">
                  {c.rank}
                </div>
              </div>

              {/* Right: Content */}
              <div className="flex-1 flex flex-col justify-between py-1">
                <div>
                  <h3 className="text-white font-bold text-lg leading-tight line-clamp-2 mb-1">
                    {c.personalInfo?.fullName}
                  </h3>
                  <div className="flex items-center gap-2 text-xs text-white/70">
                    <span className="bg-white/10 px-1.5 py-0.5 rounded font-mono border border-white/10">SBD: {c.sbd}</span>
                    <span className="truncate max-w-[100px]">{c.personalInfo?.address}</span>
                  </div>
                </div>

                <div className="mt-2 flex items-center justify-between gap-3">
                  <div className="flex flex-col">
                    <span className="text-[10px] text-white/50 uppercase tracking-wider">Bình chọn</span>
                    <span className="text-yellow-400 font-bold text-xl leading-none">{c.voteCount}</span>
                  </div>
                  <button
                    onClick={() => handleVote(c.id)}
                    disabled={voteMutation.isPending}
                    className="bg-white text-blue-600 font-bold text-sm px-4 py-2 rounded-lg shadow-sm active:scale-95 transition-transform"
                  >
                    Bình chọn
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export const Route = createFileRoute('/vote/danh-sach-thi-sinh')({
  component: ContestantList,
});
