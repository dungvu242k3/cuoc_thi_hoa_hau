import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";
import { apiRequest } from "../../lib/api";
import { ADMIN_QUERY, APPROVE_MUTATION, DASHBOARD_QUERY } from "../../lib/queries";
import { useAuthStore } from "../../store/useAuthStore";

export const Route = createFileRoute("/examiner/")({
    component: ExaminerDashboard,
    validateSearch: (search: Record<string, unknown>): { tab: "scoring" | "approval" | "contestants" } => {
        return {
            tab: (search["tab"] as "scoring" | "approval" | "contestants") || "scoring",
        };
    },
});

function ExaminerDashboard() {
    const { token } = useAuthStore();
    const queryClient = useQueryClient();
    const { tab: activeTab } = Route.useSearch();
    const [selectedContestant, setSelectedContestant] = useState<any>(null);

    // 1. Data for Scoring Tab (Public List)
    const { data: scoringData, isLoading: isLoadingScoring } = useQuery({
        queryKey: ["examinerDashboard"],
        queryFn: async () => apiRequest(DASHBOARD_QUERY, {}, token),
        enabled: !!token && (activeTab === "scoring" || activeTab === "contestants"),
    });

    // 2. Data for Approval Tab
    const { data: approvalData, isLoading: isLoadingApproval } = useQuery({
        queryKey: ["adminContestants"],
        queryFn: async () => apiRequest(ADMIN_QUERY, {}, token),
        enabled: !!token && activeTab === "approval",
    });

    // 3. Approve/Reject Mutation
    const mutation = useMutation({
        mutationFn: async ({ id, isApproved }: { id: string; isApproved: boolean }) => {
            return apiRequest(APPROVE_MUTATION, { id, isApproved }, token);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["adminContestants"] });
            queryClient.invalidateQueries({ queryKey: ["examinerDashboard"] });
            // Close modal if open
            if (selectedContestant) {
                setSelectedContestant(null);
            }
        },
        onError: (err) => {
            alert("Lỗi: " + (err as Error).message);
        },
    });

    const [searchTerm, setSearchTerm] = useState("");

    const contestants = scoringData?.publicContestants || [];
    const myScores = scoringData?.myScores || [];
    const scoredMap = new Set(myScores.map((s: any) => s.sbd));

    const removeAccents = (str: string) => {
        return str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").replace(/đ/g, "d").replace(/Đ/g, "D");
    }

    const filteredContestants = contestants.filter((c: any) => {
        const normalizedTerm = removeAccents(searchTerm.trim().toLowerCase());

        const rawName = c.personalInfo?.fullName || "";
        const normalizedName = removeAccents(rawName.toLowerCase());

        const rawSbd = c.sbd || "";
        const normalizedSbd = removeAccents(rawSbd.toLowerCase());

        return normalizedName.includes(normalizedTerm) || normalizedSbd.includes(normalizedTerm);
    });
    const pendingContestants = approvalData?.adminContestants || [];

    return (
        <div className="space-y-6">
            {/* Content: SCORING TAB (Public List) */}
            {activeTab === "scoring" && (
                <div className="space-y-6">
                    {/* Search Input */}
                    <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex items-center gap-3">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-slate-400">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
                        </svg>
                        <input
                            type="text"
                            placeholder="Tìm kiếm theo Tên hoặc SBD..."
                            className="flex-1 bg-transparent outline-none text-slate-700 placeholder:text-slate-400 font-medium"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        {searchTerm && (
                            <button onClick={() => setSearchTerm("")} className="text-slate-400 hover:text-slate-600 transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-5 h-5">
                                    <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
                                </svg>
                            </button>
                        )}
                    </div>

                    {isLoadingScoring ? (
                        <div className="text-center py-10 text-slate-500">Đang tải danh sách...</div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                            {filteredContestants.length === 0 && (
                                <div className="col-span-full text-center py-10 text-slate-400 italic">
                                    {searchTerm ? "Không tìm thấy kết quả phù hợp." : "Chưa có thí sinh nào được duyệt công khai."}
                                </div>
                            )}
                            {filteredContestants.map((c: any) => {
                                const isScored = c.sbd && scoredMap.has(c.sbd);
                                return (
                                    <Link
                                        key={c.id}
                                        to="/examiner/score/$contestantId"
                                        params={{ contestantId: c.id }}
                                        className="group relative bg-white rounded-2xl overflow-hidden border border-slate-200 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
                                    >
                                        {/* Image Area */}
                                        <div className="aspect-[3/4] overflow-hidden bg-slate-100 relative">
                                            <img
                                                src={c.portfolio?.avatarUrl || "https://placehold.co/300"}
                                                alt={c.personalInfo.fullName}
                                                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                                                onError={(e) => (e.currentTarget.src = "https://placehold.co/300")}
                                            />
                                            <div className="absolute top-3 right-3 bg-white/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-bold shadow-sm text-slate-800 border border-slate-100">
                                                SBD: {c.sbd || "---"}
                                            </div>

                                            {/* Status Overlay */}
                                            {isScored && (
                                                <div className="absolute inset-0 bg-black/10 flex items-center justify-center">
                                                    <div className="bg-emerald-500 text-white px-4 py-2 rounded-full font-bold shadow-lg transform -rotate-12 border-2 border-white">
                                                        ĐÃ CHẤM
                                                    </div>
                                                </div>
                                            )}
                                        </div>

                                        {/* Info Area */}
                                        <div className="p-4 text-center relative z-10 bg-white">
                                            <h3 className="text-lg font-bold text-slate-900 line-clamp-1 group-hover:text-purple-600 transition-colors font-serif">
                                                {c.personalInfo.fullName}
                                            </h3>
                                            <p className="text-xs text-slate-500 mt-1 uppercase tracking-wider">
                                                {c.personalInfo.job || "Thí sinh tự do"}
                                            </p>

                                            <div className={`mt-4 py-2 rounded-lg text-sm font-bold transition-colors ${isScored
                                                ? "bg-emerald-50 text-emerald-600 border border-emerald-100"
                                                : "bg-slate-900 text-white group-hover:bg-purple-600 shadow-lg group-hover:shadow-purple-200"
                                                }`}>
                                                {isScored ? "Xem lại điểm" : "Chấm điểm ngay"}
                                            </div>
                                        </div>
                                    </Link>
                                )
                            })}
                        </div>
                    )}
                </div>
            )}



            {/* Content: CONTESTANTS TAB (Read-only List) */}
            {
                activeTab === "contestants" && (
                    <div className="space-y-6">
                        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                            <div className="px-6 py-4 border-b border-slate-100 bg-slate-50">
                                <h3 className="font-semibold text-slate-800">Danh sách tất cả thí sinh</h3>
                            </div>
                            <div className="divide-y divide-slate-100">
                                {contestants.map((c: any) => (
                                    <div key={c.id} className="p-4 flex items-center gap-4 hover:bg-slate-50 transition-colors">
                                        <img
                                            src={c.portfolio?.avatarUrl || "https://placehold.co/150"}
                                            alt={c.personalInfo.fullName}
                                            className="w-12 h-12 rounded-full object-cover border border-slate-200"
                                            onError={(e) => (e.currentTarget.src = "https://placehold.co/150")}
                                        />
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-baseline gap-2">
                                                <h4 className="font-bold text-slate-900 truncate">{c.personalInfo.fullName}</h4>
                                                <span className="text-xs font-mono text-slate-500 bg-slate-100 px-1.5 rounded">SBD: {c.sbd || "---"}</span>
                                            </div>
                                            <p className="text-sm text-slate-500 truncate">{c.personalInfo.job || "Thí sinh tự do"}</p>
                                        </div>
                                        <button
                                            onClick={() => setSelectedContestant(c)}
                                            className="px-3 py-1.5 text-xs font-bold text-slate-600 bg-white border border-slate-200 rounded-lg hover:border-purple-300 hover:text-purple-600 transition-all"
                                        >
                                            Xem hồ sơ
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                )
            }

            {/* Content: APPROVAL TAB */}
            {
                activeTab === "approval" && (
                    <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                        <div className="px-6 py-4 border-b border-slate-100 bg-slate-50 flex justify-between items-center">
                            <h3 className="font-semibold text-slate-800">Danh sách chờ duyệt</h3>
                            <span className="text-xs font-medium px-2 py-1 bg-amber-100 text-amber-700 rounded-full">
                                Pending: {pendingContestants.length}
                            </span>
                        </div>

                        {isLoadingApproval ? (
                            <div className="p-10 text-center text-slate-500">Đang tải hồ sơ chờ...</div>
                        ) : pendingContestants.length === 0 ? (
                            <div className="p-10 text-center text-slate-400 italic">
                                Không có hồ sơ nào đang chờ duyệt.
                            </div>
                        ) : (
                            <div className="divide-y divide-slate-100">
                                {pendingContestants.map((c: any) => (
                                    <div key={c.id} className="p-6 flex flex-col sm:flex-row gap-6 hover:bg-slate-50 transition-colors">
                                        <div className="shrink-0">
                                            <img
                                                src={c.portfolio?.avatarUrl || "https://placehold.co/150"}
                                                className="w-24 h-24 rounded-lg object-cover shadow-sm bg-slate-200"
                                                alt="Avatar"
                                            />
                                        </div>
                                        <div className="flex-1 space-y-2">
                                            <div className="flex justify-between items-start">
                                                <div>
                                                    <h4 className="text-lg font-bold text-slate-900">{c.personalInfo.fullName}</h4>
                                                    <p className="text-sm text-slate-500">{c.personalInfo.job} • {new Date(c.createdAt).toLocaleDateString()}</p>
                                                </div>
                                                <div className="flex gap-2">
                                                    <button
                                                        onClick={() => setSelectedContestant(c)}
                                                        className="px-3 py-1 text-sm font-medium text-purple-600 bg-purple-50 hover:bg-purple-100 rounded-md border border-purple-200 transition-colors"
                                                    >
                                                        Xem chi tiết
                                                    </button>
                                                    <button
                                                        onClick={() => mutation.mutate({ id: c.id, isApproved: false })}
                                                        disabled={mutation.isPending}
                                                        className="px-3 py-1 text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-md border border-red-200 transition-colors"
                                                    >
                                                        Từ chối
                                                    </button>
                                                    <button
                                                        onClick={() => mutation.mutate({ id: c.id, isApproved: true })}
                                                        disabled={mutation.isPending}
                                                        className="px-4 py-1 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-700 rounded-md shadow-sm transition-all"
                                                    >
                                                        {mutation.isPending ? "..." : "Duyệt"}
                                                    </button>
                                                </div>
                                            </div>

                                            <div className="bg-slate-50 p-3 rounded-md border border-slate-100 text-sm text-slate-600">
                                                <p className="font-semibold text-xs text-slate-400 uppercase mb-1">Giới thiệu bản thân</p>
                                                <p className="line-clamp-2">{c.portfolio?.introduction || "Không có mô tả"}</p>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )
            }



            {/* Modal */}
            {
                selectedContestant && (
                    <ContestantDetailModal
                        contestant={selectedContestant}
                        onClose={() => setSelectedContestant(null)}
                        onApprove={() => mutation.mutate({ id: selectedContestant.id, isApproved: true })}
                        onReject={() => mutation.mutate({ id: selectedContestant.id, isApproved: false })}
                        isProcessing={mutation.isPending}
                    />
                )
            }
        </div>
    );
}

function ContestantDetailModal({
    contestant,
    onClose,
    onApprove,
    onReject,
    isProcessing
}: {
    contestant: any;
    onClose: () => void;
    onApprove: () => void;
    onReject: () => void;
    isProcessing: boolean;
}) {
    if (!contestant) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm overflow-y-auto">
            <div className="bg-white rounded-2xl shadow-xl w-full max-w-4xl max-h-[90vh] overflow-y-auto flex flex-col animate-in fade-in zoom-in-95 duration-200">
                {/* Header */}
                <div className="p-6 border-b border-slate-100 flex justify-between items-center sticky top-0 bg-white z-10">
                    <div>
                        <h3 className="text-xl font-bold text-slate-900">{contestant.personalInfo.fullName}</h3>
                        <p className="text-sm text-slate-500">SBD: {contestant.sbd || "Chưa cấp"} • {new Date(contestant.createdAt).toLocaleDateString()}</p>
                    </div>
                    <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-full text-slate-400 hover:text-slate-600 transition-colors">
                        ✕
                    </button>
                </div>

                {/* Body */}
                <div className="p-6 space-y-8 flex-1 overflow-y-auto">
                    {/* 1. Basic Info & Portfolio */}
                    <div className="flex flex-col md:flex-row gap-8">
                        <div className="shrink-0">
                            <img
                                src={contestant.portfolio?.avatarUrl || "https://placehold.co/300"}
                                alt="Avatar"
                                className="w-64 h-80 object-cover rounded-xl shadow-md bg-slate-100"
                            />
                        </div>
                        <div className="flex-1 space-y-6">
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Ngày sinh</span>
                                    <span className="text-slate-700">{new Date(contestant.personalInfo.dob).toLocaleDateString()}</span>
                                </div>
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Giới tính</span>
                                    <span className="text-slate-700">{contestant.personalInfo.gender || "N/A"}</span>
                                </div>
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Quốc tịch</span>
                                    <span className="text-slate-700">{contestant.personalInfo.nationality}</span>
                                </div>
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Nghề nghiệp</span>
                                    <span className="text-slate-700">{contestant.personalInfo.job}</span>
                                </div>
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">SĐT</span>
                                    <span className="text-slate-700">{contestant.personalInfo.phone}</span>
                                </div>
                                <div>
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Email</span>
                                    <span className="text-slate-700">{contestant.personalInfo.email}</span>
                                </div>

                                <div className="col-span-2">
                                    <span className="block text-xs font-semibold text-slate-400 uppercase">Địa chỉ</span>
                                    <span className="text-slate-700">{contestant.personalInfo.address}</span>
                                </div>
                            </div>

                            <div className="bg-purple-50 p-4 rounded-xl border border-purple-100">
                                <h4 className="font-bold text-purple-900 mb-2">Chỉ số hình thể</h4>
                                <div className="grid grid-cols-3 gap-4 text-center">
                                    <div>
                                        <span className="block text-2xl font-bold text-purple-700">{contestant.physicalInfo?.height || 0}</span>
                                        <span className="text-xs text-purple-600">Chiều cao (cm)</span>
                                    </div>
                                    <div>
                                        <span className="block text-2xl font-bold text-purple-700">{contestant.physicalInfo?.weight || 0}</span>
                                        <span className="text-xs text-purple-600">Cân nặng (kg)</span>
                                    </div>
                                    <div>
                                        <span className="block text-2xl font-bold text-purple-700">{contestant.physicalInfo?.measurements || "-"}</span>
                                        <span className="text-xs text-purple-600">Số đo 3 vòng</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* 2. Introduction */}
                    <div>
                        <h4 className="font-bold text-slate-800 mb-2 text-lg">Giới thiệu bản thân</h4>
                        <p className="text-slate-600 leading-relaxed bg-slate-50 p-4 rounded-xl border border-slate-100">
                            {contestant.portfolio?.introduction || "Chưa cập nhật giới thiệu."}
                        </p>
                    </div>

                    {/* 3. Skill & Edu */}
                    <div className="grid md:grid-cols-2 gap-6">
                        <div>
                            <h4 className="font-bold text-slate-800 mb-2">Học vấn & Ngôn ngữ</h4>
                            <ul className="list-disc list-inside text-slate-600 space-y-1">
                                <li>
                                    Trình độ: <span className="font-medium text-slate-800">
                                        {(() => {
                                            const level = contestant.skillEducation?.educationLevel?.toLowerCase();
                                            const map: Record<string, string> = {
                                                "university": "Đại học",
                                                "college": "Cao đẳng",
                                                "highschool": "Trung học phổ thông",
                                                "master": "Thạc sĩ",
                                                "doctor": "Tiến sĩ",
                                                "intermediate": "Trung cấp",
                                                "other": "Khác"
                                            };
                                            return map[level] || contestant.skillEducation?.educationLevel || "Chưa cập nhật";
                                        })()}
                                    </span>
                                </li>
                                <li>Ngoại ngữ: {contestant.skillEducation?.languages?.join(", ") || "Không"}</li>
                            </ul>
                        </div>
                        <div>
                            <h4 className="font-bold text-slate-800 mb-2">Kỹ năng năng khiếu</h4>
                            <div className="flex flex-wrap gap-2">
                                {contestant.skillEducation?.skills?.length > 0 ? (
                                    contestant.skillEducation.skills.map((s: string, idx: number) => (
                                        <span key={idx} className="px-3 py-1 bg-emerald-50 text-emerald-700 rounded-full text-sm font-medium border border-emerald-100">
                                            {s}
                                        </span>
                                    ))
                                ) : (
                                    <span className="text-slate-400 italic">Chưa cập nhật kỹ năng</span>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* 4. Gallery */}
                    {contestant.portfolio?.galleryUrls?.length > 0 && (
                        <div>
                            <h4 className="font-bold text-slate-800 mb-3">Thư viện ảnh</h4>
                            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                                {contestant.portfolio.galleryUrls
                                    .filter((url: string) => url && url.trim() !== "")
                                    .map((url: string, idx: number) => (
                                        <img
                                            key={idx}
                                            src={url}
                                            alt={`Gallery ${idx}`}
                                            className="w-full aspect-[3/4] object-cover rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer border border-slate-100"
                                            onClick={() => window.open(url, '_blank')}
                                            onError={(e) => e.currentTarget.style.display = 'none'}
                                        />
                                    ))}
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer Actions */}
                {/* Footer Actions */}
                <div className="p-6 border-t border-slate-100 bg-slate-50 sticky bottom-0 flex justify-end gap-3 rounded-b-2xl">
                    <button
                        onClick={onClose}
                        className="px-5 py-2.5 font-medium text-slate-600 hover:bg-white hover:shadow-sm rounded-lg border border-transparent hover:border-slate-200 transition-all"
                    >
                        Đóng
                    </button>

                    {/* Only show Approve/Reject if NOT public/approved */}
                    {!contestant.isPublic && contestant.status !== 'approved' && (
                        <>
                            <button
                                onClick={onReject}
                                disabled={isProcessing}
                                className="px-5 py-2.5 font-medium text-red-600 bg-white border border-red-200 hover:bg-red-50 rounded-lg shadow-sm transition-all disabled:opacity-50"
                            >
                                Từ chối hồ sơ
                            </button>
                            <button
                                onClick={onApprove}
                                disabled={isProcessing}
                                className="px-6 py-2.5 font-bold text-white bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700 rounded-lg shadow-md hover:shadow-lg transition-all disabled:opacity-50"
                            >
                                {isProcessing ? "Đang xử lý..." : "DUYỆT HỒ SƠ"}
                            </button>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}
