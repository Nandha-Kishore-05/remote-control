export default function StatsCard({ title, value, color }) {
  const bgColors = {
    blue: "bg-blue-50 text-blue-800",
    green: "bg-emerald-50 text-emerald-800",
    red: "bg-rose-50 text-rose-800",
    amber: "bg-amber-50 text-amber-800",
    indigo: "bg-indigo-50 text-indigo-800",
  };

  const classes = bgColors[color] || "bg-gray-50 text-gray-800";

  return (
    <div className={`rounded-lg shadow p-4 ${classes}`}>
      <h3 className="text-sm font-medium">{title}</h3>
      <p className="mt-1 text-3xl font-semibold">{value}</p>
    </div>
  );
}
