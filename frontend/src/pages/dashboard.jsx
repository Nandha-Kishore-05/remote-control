import { useState } from "react";
import StatsCard from "../components/StatsCard";
import PCList from "../components/pclist";
import PCModal from "../components/pcmodal";
import ConfirmDeleteModal from "../components/confirmDeleteModal.";

const initialData = [
  { id: 1, name: "PC-01", ip: "192.168.1.10", mac: "AA:BB:CC:DD:EE:FF", status: "Turn On", seb: true },
  { id: 2, name: "PC-02", ip: "192.168.1.11", mac: "AA:BB:CC:DD:EE:00", status: "Turn Off", seb: false },
  { id: 3, name: "PC-03", ip: "192.168.1.12", mac: "AA:BB:CC:DD:EE:11", status: "Turn On", seb: false },
];

export default function Dashboard() {
  const [pcs, setPcs] = useState(initialData);
  const [selectedPC, setSelectedPC] = useState(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  const handleAdd = () => {
    setSelectedPC(null);
    setIsFormOpen(true);
  };

  const handleEdit = (pc) => {
    setSelectedPC(pc);
    setIsFormOpen(true);
  };

  const handleDelete = (pc) => {
    setSelectedPC(pc);
    setIsDeleteOpen(true);
  };

  const handleSave = (pc) => {
    if (pc.id) {
      setPcs((prev) =>
        prev.map((item) => (item.id === pc.id ? pc : item))
      );
    } else {
      const newPC = { ...pc, id: Date.now(), status: "Offline", seb: false };
      setPcs((prev) => [...prev, newPC]);
    }
    setIsFormOpen(false);
  };

  const handleConfirmDelete = () => {
    setPcs((prev) => prev.filter((item) => item.id !== selectedPC.id));
    setIsDeleteOpen(false);
  };

  // Stats
  const total = pcs.length;
  const online = pcs.filter((p) => p.status === "Online").length;
  const offline = pcs.filter((p) => p.status === "Offline").length;
  const sebRunning = pcs.filter((p) => p.seb).length;

  return (
    <div className="p-6 bg-[#f4f6f8] min-h-screen">

      <h1 className="text-3xl font-bold text-[#444] mb-6">💻 Remote Control</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <StatsCard title="Total Systems" value={12} color="indigo" />
        <StatsCard title="Online" value={7} color="green" />
        <StatsCard title="Offline" value={5} color="red" />
        <StatsCard title="SEB running" value={7} color="green" />

      </div>


      <PCList pcs={pcs} onEdit={handleEdit} onDelete={handleDelete} />


      <PCModal
        isOpen={isFormOpen}
        onClose={() => setIsFormOpen(false)}
        onSave={handleSave}
        pc={selectedPC}
      />

      <ConfirmDeleteModal
        isOpen={isDeleteOpen}
        onClose={() => setIsDeleteOpen(false)}
        onConfirm={handleConfirmDelete}
        pc={selectedPC}
      />
    </div>
  );
}
