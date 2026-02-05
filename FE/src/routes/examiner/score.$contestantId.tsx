import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { apiRequest } from "../../lib/api";
import { APP_CONFIG, SCORING_CRITERIA } from "../../lib/constants";
import { GET_CONTESTANT, SUBMIT_SCORE } from "../../lib/queries";
import { useAuthStore } from "../../store/useAuthStore";

export const Route = createFileRoute("/examiner/score/$contestantId")({
    component: ScoringPage,
});

function ScoringPage() {
    const { contestantId } = Route.useParams();
    const { token } = useAuthStore();
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    // Dynamic State Initialization
    const [scores, setScores] = useState<Record<string, string>>({});
    const [comment, setComment] = useState("");

    // Fetch Contestant Data
    const { data, isLoading } = useQuery({
        queryKey: ["contestant", contestantId],
        queryFn: async () => apiRequest(GET_CONTESTANT, { id: contestantId }, token),
        enabled: !!token
    });

    // Submit Mutation
    const mutation = useMutation({
        mutationFn: async (payload: any) => apiRequest(SUBMIT_SCORE, { input: payload }, token),
        onSuccess: () => {
            alert("Đã lưu điểm thành công!");
            queryClient.invalidateQueries({ queryKey: ["examinerDashboard"] });
            navigate({ to: "/examiner", search: { tab: "scoring" } });
        },
        onError: (err) => {
            alert("Lỗi khi lưu điểm: " + err);
        }
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        const criteriaScores: Record<string, number> = {};
        SCORING_CRITERIA.forEach(c => {
            criteriaScores[c.key] = parseFloat(scores[c.key] ?? "") || 0;
        });

        const payload = {
            contestantId,
            sbd: data?.publicContestantDetail?.sbd || "UNKNOWN",
            criteriaScores,
            comment
        };
        mutation.mutate(payload);
    };

    if (isLoading) return <div className="p-10 text-center">Đang tải hồ sơ...</div>;

    const c = data?.publicContestantDetail;
    if (!c) return <div className="p-10 text-center text-red-500">Không tìm thấy thí sinh</div>;

    // Dynamic Total Calculation
    const totalScore = SCORING_CRITERIA.reduce((sum, item) => sum + (Number(scores[item.key]) || 0), 0);

    // Helper to extract YouTube ID if needed
    const getYoutubeEmbed = (url: string) => {
        if (!url) return null;
        const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
        const match = url.match(regExp);
        return (match && match[2] && match[2].length === 11) ? `https://www.youtube.com/embed/${match[2]}` : null;
    };
    const embedUrl = getYoutubeEmbed(c.portfolio?.videoUrl || "");


    return (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[calc(100vh-100px)] lg:h-[calc(100vh-100px)]">
            {/* LEFT COLUMN: Contestant Profile */}
            <div className="lg:col-span-5 xl:col-span-4 space-y-6 lg:overflow-y-auto lg:pr-2 custom-scrollbar lg:pb-10">
                {/* Back Link */}
                <Link to="/examiner" search={{ tab: "scoring" }} className="inline-flex items-center text-sm text-slate-500 hover:text-slate-800 transition-colors">
                    &larr; Quay lại danh sách
                </Link>

                {/* 1. Main Info Card */}
                <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
                    <div className="aspect-[3/4] w-full relative">
                        <img
                            src={c.portfolio?.avatarUrl || "https://placehold.co/300"}
                            alt={c.personalInfo.fullName}
                            className="absolute inset-0 w-full h-full object-cover"
                            onError={(e) => (e.currentTarget.src = "https://placehold.co/300")}
                        />
                        <div className="absolute top-4 left-4 bg-white/90 backdrop-blur px-3 py-1 rounded-full text-xs font-bold shadow-sm border border-slate-200">
                            SBD: {c.sbd || "N/A"}
                        </div>
                    </div>
                    <div className="p-5">
                        <h2 className="text-xl font-bold text-slate-900">{c.personalInfo.fullName}</h2>
                        <p className="text-sm text-slate-500 mb-4">{c.personalInfo.job}</p>

                        <div className="space-y-3 text-sm">
                            <div className="flex justify-between border-b border-slate-50 pb-2">
                                <span className="text-slate-500">Ngày sinh</span>
                                <span className="font-medium text-slate-700">{c.personalInfo.dob ? new Date(c.personalInfo.dob).toLocaleDateString() : "N/A"}</span>
                            </div>
                            <div className="flex justify-between border-b border-slate-50 pb-2">
                                <span className="text-slate-500">Quê quán</span>
                                <span className="font-medium text-slate-700">{c.personalInfo.hometown}</span>
                            </div>
                            <div className="flex justify-between border-b border-slate-50 pb-2">
                                <span className="text-slate-500">Hình thể</span>
                                <span className="font-medium text-slate-700">{c.physicalInfo?.height}cm - {c.physicalInfo?.weight}kg</span>
                            </div>
                            <div className="flex justify-between border-b border-slate-50 pb-2">
                                <span className="text-slate-500">Số đo 3 vòng</span>
                                <span className="font-medium text-slate-700">{c.physicalInfo?.measurements}</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* 2. Video Player */}
                {c.portfolio?.videoUrl && (
                    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-4">
                        <h3 className="font-bold text-slate-800 mb-3 flex items-center gap-2">
                            Video giới thiệu
                        </h3>
                        {embedUrl ? (
                            <iframe
                                src={embedUrl}
                                className="w-full aspect-video rounded-lg bg-black"
                                title="Contestant Video"
                                allowFullScreen
                            />
                        ) : (
                            <div className="p-4 bg-slate-50 rounded-lg text-sm text-slate-500 truncate">
                                <a href={c.portfolio.videoUrl} target="_blank" rel="noreferrer" className="text-blue-600 hover:underline">
                                    Xem video tại đây (Link ngoài)
                                </a>
                            </div>
                        )}
                    </div>
                )}

                {/* 3. Introduction & Skills */}
                <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-5 space-y-4">
                    <div>
                        <h3 className="font-bold text-slate-800 mb-2">Giới thiệu bản thân</h3>
                        <p className="text-sm text-slate-600 bg-slate-50 p-3 rounded-lg border border-slate-100 leading-relaxed whitespace-pre-line">
                            {c.portfolio?.introduction || "Chưa cập nhật giới thiệu."}
                        </p>
                    </div>
                    <div>
                        <h3 className="font-bold text-slate-800 mb-2">Kỹ năng & Học vấn</h3>
                        <ul className="text-sm text-slate-600 space-y-1 ml-4 list-disc">
                            <li><span className="font-medium">Trình độ:</span> {c.skillEducation?.educationLevel || "N/A"}</li>
                            <li><span className="font-medium">Ngoại ngữ:</span> {c.skillEducation?.languages?.join(", ") || "Không"}</li>
                            <li><span className="font-medium">Năng khiếu:</span> {c.skillEducation?.skills?.join(", ") || "Không"}</li>
                        </ul>
                    </div>
                </div>

                {/* 4. Gallery Preview */}
                {c.portfolio?.galleryUrls && c.portfolio.galleryUrls.length > 0 && (
                    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-4">
                        <h3 className="font-bold text-slate-800 mb-3">Thư viện ảnh</h3>
                        <div className="grid grid-cols-3 gap-2">
                            {c.portfolio.galleryUrls
                                .filter((url: string) => url && url.trim() !== "")
                                .map((img: string, i: number) => (
                                    <img
                                        key={i}
                                        src={img}
                                        className="w-full aspect-square object-cover rounded-lg border border-slate-200 shadow-sm cursor-pointer hover:opacity-90 active:scale-95 transition-all"
                                        onClick={() => window.open(img, '_blank')}
                                        onError={(e) => e.currentTarget.style.display = 'none'}
                                    />
                                ))}
                        </div>
                    </div>
                )}
            </div>

            {/* RIGHT COLUMN: Scoring Form */}
            <div className="lg:col-span-7 xl:col-span-8 flex flex-col h-full lg:overflow-hidden pb-6">
                <div className="bg-white rounded-2xl shadow-lg border border-slate-200 flex-1 flex flex-col overflow-hidden">
                    <div className="p-5 border-b border-slate-100 bg-slate-50/50">
                        <h3 className="font-bold text-lg text-slate-800 flex items-center gap-2">
                            Phiếu chấm điểm
                        </h3>
                        <p className="text-xs text-slate-500 mt-1">Vui lòng chấm điểm công tâm và chính xác theo thang điểm 10.</p>
                    </div>

                    <div className="p-6 md:p-8 lg:overflow-y-auto flex-1 custom-scrollbar">
                        <form id="score-form" onSubmit={handleSubmit} className="max-w-3xl mx-auto space-y-8">

                            {/* Scoring Grid - DYNAMIC RENDER */}
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                                {SCORING_CRITERIA.map((criterion) => (
                                    <ScoreInput
                                        key={criterion.key}
                                        label={criterion.label}
                                        value={scores[criterion.key] || ""}
                                        onChange={(v) => setScores({ ...scores, [criterion.key]: v })}
                                        desc={criterion.desc}
                                    />
                                ))}
                            </div>

                            {/* Comment Section */}
                            <div className="pt-6 border-t border-slate-100">
                                <label className="block text-sm font-bold text-slate-700 mb-2">
                                    Ghi chú / Nhận xét thêm
                                </label>
                                <textarea
                                    value={comment}
                                    onChange={(e) => setComment(e.target.value)}
                                    rows={4}
                                    className="block w-full rounded-xl border-slate-200 bg-slate-50 focus:bg-white focus:border-purple-500 focus:ring-purple-500 transition-all text-sm p-4 placeholder:text-slate-400"
                                    placeholder="Nhập nhận xét chi tiết về thí sinh này..."
                                />
                            </div>

                        </form>
                    </div>

                    {/* Sticky Footer */}
                    <div className="p-4 border-t border-slate-200 bg-white flex justify-end gap-3 items-center sticky bottom-0 z-10 lg:static">
                        <span className="text-sm font-medium text-slate-500 mr-2">
                            Tổng điểm: <span className="text-purple-600 font-bold text-lg">
                                {totalScore.toFixed(1)}
                            </span>
                        </span>
                        <button
                            form="score-form"
                            type="submit"
                            disabled={mutation.isPending}
                            className="px-8 py-2.5 rounded-lg bg-gradient-to-r from-purple-600 to-rose-500 text-white font-bold shadow-lg shadow-purple-200 hover:shadow-xl hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50"
                        >
                            {mutation.isPending ? "Đang lưu..." : "Xác nhận chấm"}
                        </button>
                    </div>
                </div>
            </div>
        </div >
    );
}

function ScoreInput({ label, value, onChange, desc }: { label: string, value: string, onChange: (v: string) => void, desc: string }) {
    return (
        <div className="group">
            <div className="flex justify-between mb-1.5">
                <label className="block text-sm font-bold text-slate-700 group-focus-within:text-purple-700 transition-colors">
                    {label}
                </label>
                <span className="text-xs font-mono text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded">Max: {APP_CONFIG.MAX_SCORE}</span>
            </div>
            <input
                type="number"
                min={APP_CONFIG.MIN_SCORE}
                max={APP_CONFIG.MAX_SCORE}
                step={APP_CONFIG.STEP}
                value={value}
                onChange={(e) => {
                    const val = Math.min(APP_CONFIG.MAX_SCORE, Math.max(APP_CONFIG.MIN_SCORE, Number(e.target.value)));
                    onChange(val.toString());
                }}
                className="block w-full rounded-lg border-slate-200 focus:border-purple-500 focus:ring-purple-500 text-lg font-semibold py-3 px-4 shadow-sm transition-all"
                placeholder="0.0"
            />
            <p className="mt-1.5 text-xs text-slate-500 pl-1">{desc}</p>
        </div>
    )
}
