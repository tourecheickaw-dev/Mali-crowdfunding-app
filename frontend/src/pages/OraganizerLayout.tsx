import { NavLink, Outlet } from 'react-router-dom'

const tab = ({ isActive }: { isActive: boolean }) =>
  `rounded-full px-5 py-2 font-semibold ${isActive ? 'bg-[#111827] text-white' : 'bg-gray-100 text-gray-700'}`

export default function OrganizerLayout() {
  return (
    <div className="mx-auto max-w-2xl py-6">
      <div className="flex gap-2">
        <NavLink to="/mon-espace" end className={tab}>Mes cagnottes</NavLink>
        <NavLink to="/mon-espace/paiement" className={tab}>Paiement</NavLink>
      </div>
      <Outlet />
    </div>
  )
}