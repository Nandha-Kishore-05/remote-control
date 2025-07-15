import { useState, useEffect } from "react";

export default function PCModal({ isOpen, onClose, onSave, pc }) {
  const [form, setForm] = useState({ name: "", ip: "", mac: "" });

  useEffect(() => {
    if (pc) setForm(pc);
    else setForm({ name: "", ip: "", mac: "" });
  }, [pc]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg p-6 animate-fadeIn">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">
          {pc ? "Edit PC" : "Add PC"}
        </h2>
        <div className="space-y-4">
          <input
            className="w-full border border-gray-300 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 px-4 py-2 rounded-lg transition placeholder-gray-400"
            placeholder="Name"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
          />
          <input
            className="w-full border border-gray-300 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 px-4 py-2 rounded-lg transition placeholder-gray-400"
            placeholder="IP Address"
            value={form.ip}
            onChange={(e) => setForm({ ...form, ip: e.target.value })}
          />
          <input
            className="w-full border border-gray-300 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 px-4 py-2 rounded-lg transition placeholder-gray-400"
            placeholder="MAC Address"
            value={form.mac}
            onChange={(e) => setForm({ ...form, mac: e.target.value })}
          />
        </div>
        <div className="flex justify-end mt-6 space-x-3">
          <button
            className="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-lg transition"
            onClick={onClose}
          >
            Cancel
          </button>
          <button
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
            onClick={() => onSave(form)}
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
}
