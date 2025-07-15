import { useState, useMemo } from "react";
import EditOutlinedIcon from "@mui/icons-material/EditOutlined";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import KeyboardArrowLeftIcon from "@mui/icons-material/KeyboardArrowLeft";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import PCModal from "./pcmodal";


export default function PCList({ pcs, onEdit, onDelete }) {
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [sebFilter, setSebFilter] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedIds, setSelectedIds] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingPC, setEditingPC] = useState(null);

  const bgColors = {
    blue: "bg-blue-50 text-blue-800",
    green: "bg-emerald-50 text-emerald-800",
    red: "bg-rose-50 text-rose-800",
    amber: "bg-amber-50 text-amber-800",
    indigo: "bg-indigo-50 text-indigo-800",
  };

  const filteredPCs = useMemo(() => {
    return pcs.filter((pc) => {
      const matchesSearch =
        pc.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        pc.ip.toLowerCase().includes(searchTerm.toLowerCase()) ||
        pc.mac.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesStatus = statusFilter ? pc.status === statusFilter : true;
      const matchesSeb =
        sebFilter === ""
          ? true
          : sebFilter === "Running"
          ? pc.seb
          : !pc.seb;
      return matchesSearch && matchesStatus && matchesSeb;
    });
  }, [pcs, searchTerm, statusFilter, sebFilter]);

  const totalPages = Math.ceil(filteredPCs.length / rowsPerPage);
  const paginatedPCs = useMemo(() => {
    const startIndex = (currentPage - 1) * rowsPerPage;
    return filteredPCs.slice(startIndex, startIndex + rowsPerPage);
  }, [filteredPCs, currentPage, rowsPerPage]);

  const allVisibleSelected = paginatedPCs.every((pc) =>
    selectedIds.includes(pc.id)
  );

  const anySelected = selectedIds.length > 0;

  // Save handler
  const handleSave = (formData) => {
    if (editingPC) {
      // update existing
      onEdit({ ...editingPC, ...formData });
    } else {
      // create new
      onEdit(formData);
    }
    setModalOpen(false);
  };

  return (
    <div className="bg-white p-4 rounded-xl shadow ring-1 ring-gray-200">
      {/* Filters */}
      <div className="flex flex-wrap items-center justify-between mb-4 gap-2">
        <input
          type="text"
          placeholder="Search by Name, IP, or MAC..."
          value={searchTerm}
          onChange={(e) => {
            setSearchTerm(e.target.value);
            setCurrentPage(1);
            setSelectedIds([]);
          }}
          className="border border-gray-300 rounded-lg px-3 py-2 w-full md:w-64 focus:outline-none focus:ring focus:border-blue-300 text-sm"
        />
        <div className="flex gap-2 flex-wrap">
          <select
            value={statusFilter}
            onChange={(e) => {
              setStatusFilter(e.target.value);
              setCurrentPage(1);
              setSelectedIds([]);
            }}
            className="border border-gray-300 rounded-lg px-3 py-2 text-sm"
          >
            <option value="">All Status</option>
            <option value="Online">Online</option>
            <option value="Offline">Offline</option>
          </select>
          <select
            value={sebFilter}
            onChange={(e) => {
              setSebFilter(e.target.value);
              setCurrentPage(1);
              setSelectedIds([]);
            }}
            className="border border-gray-300 rounded-lg px-3 py-2 text-sm"
          >
            <option value="">All SEB</option>
            <option value="Running">Running</option>
            <option value="Stopped">Stopped</option>
          </select>
          <button
            className="px-4 py-2 bg-blue-500 rounded-lg text-white cursor-pointer transition"
            onClick={() => {
              setEditingPC(null);
              setModalOpen(true);
            }}
          >
            + Add System
          </button>
        </div>
      </div>

      {/* Bulk Action */}
      {anySelected && (
        <div className="mb-4 flex flex-wrap items-center gap-3 bg-white border border-gray-300 rounded-xl p-4 shadow">
          <span className="text-base font-semibold text-gray-800">
            {selectedIds.length} selected
          </span>
          {/* Buttons for bulk actions */}
          <button
            className={`px-4 py-2 rounded-lg text-[15px] font-medium ${bgColors.green} hover:bg-emerald-100 transition cursor-pointer`}
          >
            Power On
          </button>
          <button
            className={`px-4 py-2 rounded-lg text-[15px] font-medium ${bgColors.red} hover:bg-rose-100 transition cursor-pointer`}
          >
            Shutdown
          </button>
        </div>
      )}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="min-w-full text-sm">
          <thead className="bg-gray-100 text-gray-800">
            <tr>
              <th className="px-2 py-4 text-center rounded-tl-xl">
                <input
                  type="checkbox"
                  checked={allVisibleSelected}
                  onChange={(e) => {
                    if (e.target.checked) {
                      const newIds = paginatedPCs.map((pc) => pc.id);
                      setSelectedIds((prev) => [...new Set([...prev, ...newIds])]);
                    } else {
                      const newIds = paginatedPCs.map((pc) => pc.id);
                      setSelectedIds((prev) =>
                        prev.filter((id) => !newIds.includes(id))
                      );
                    }
                  }}
                />
              </th>
              <th className="px-6 py-4 text-left">Name</th>
              <th className="px-6 py-4 text-center">IP</th>
              <th className="px-6 py-4 text-center">MAC</th>
              <th className="px-6 py-4 text-center">Status</th>
              <th className="px-6 py-4 text-center">SEB</th>
              <th className="px-6 py-4 text-center rounded-tr-xl">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {paginatedPCs.map((pc) => (
              <tr key={pc.id}>
                <td className="px-3 py-4 text-center">
                  <input
                    type="checkbox"
                    checked={selectedIds.includes(pc.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedIds((prev) => [...prev, pc.id]);
                      } else {
                        setSelectedIds((prev) =>
                          prev.filter((id) => id !== pc.id)
                        );
                      }
                    }}
                  />
                </td>
                <td className="px-6 py-4">{pc.name}</td>
                <td className="px-6 py-4 text-center">{pc.ip}</td>
                <td className="px-6 py-4 text-center">{pc.mac}</td>
                <td className="px-6 py-4 text-center">
                  <span
                    className={`inline-flex px-3 py-1 rounded-full ${
                      pc.status === "Online"
                        ? bgColors.green
                        : bgColors.red
                    }`}
                  >
                    {pc.status}
                  </span>
                </td>
                <td className="px-6 py-4 text-center">
                  {pc.seb ? (
                    <span
                      className={`inline-flex px-3 py-1 rounded-full ${bgColors.amber}`}
                    >
                      Running
                    </span>
                  ) : (
                    <span className="text-gray-500">Stopped</span>
                  )}
                </td>
                <td className="px-6 py-4 text-center">
                  <div className="flex gap-2 justify-center">
                    <button
                      onClick={() => {
                        setEditingPC(pc);
                        setModalOpen(true);
                      }}
                      className="w-10 h-10 flex items-center justify-center rounded-lg border text-gray-700 hover:bg-gray-100"
                    >
                      <EditOutlinedIcon fontSize="small" />
                    </button>
                    <button
                      onClick={() => onDelete(pc)}
                      className="w-10 h-10 flex items-center justify-center rounded-lg border text-gray-700 hover:bg-gray-100"
                    >
                      <DeleteOutlineIcon fontSize="small" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
            {paginatedPCs.length === 0 && (
              <tr>
                <td colSpan="7" className="text-center py-6 text-gray-500">
                  No records found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between  flex-wrap gap-2 bg-gray-100 px-4 py-2 rounded-b-xl">
        <p className="text-sm">
          Page {currentPage} of {totalPages || 1}
        </p>
        <div className="flex gap-2">
          <select
            value={rowsPerPage}
            onChange={(e) => {
              setRowsPerPage(Number(e.target.value));
              setCurrentPage(1);
            }}
            className="border rounded px-2 py-1 text-sm"
          >
            <option value={5}>5 / page</option>
            <option value={10}>10 / page</option>
            <option value={20}>20 / page</option>
          </select>
          <button
            onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
            disabled={currentPage === 1}
            className="border rounded px-2 py-1 disabled:opacity-50"
          >
            <KeyboardArrowLeftIcon fontSize="small" />
          </button>
          <button
            onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
            disabled={currentPage === totalPages || totalPages === 0}
            className="border rounded px-2 py-1 disabled:opacity-50"
          >
            <KeyboardArrowRightIcon fontSize="small" />
          </button>
        </div>
      </div>

      {/* Modal */}
      <PCModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        onSave={handleSave}
        pc={editingPC}
      />
    </div>
  );
}
