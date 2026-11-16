import { EyeOff } from 'lucide-react'
import { Link } from 'react-router-dom'
import { fcfa } from '../api'
import { UEMOA_COUNTRIES } from '../constants'
import type { Campaign } from '../types'

export default function CampaignCard({ c, big = false }: { c: Campaign; big?: boolean }) {
  const goal = c.goal_amount ?? 0
  const pct = goal > 0 ? Math.min(100, Math.round((c.current_amount / goal) * 100)) : 0
  const country = UEMOA_COUNTRIES.find((x) => x.code === c.country)?.name

  return (
    <Link
      to={`/c/${c.url_slug}`}
      className="group block overflow-hidden rounded-3xl bg-white shadow-sm ring-1 ring-gray-200 transition hover:-translate-y-1 hover:shadow-lg"
    >
      <div className="relative aspect-4/3 overflow-hidden bg-brand-600rand-100">
        {c.cover_url ? (
          <img
            src={c.cover_url}
            alt=""
            loading="lazy"
            className={`h-full w-full object-cover transition duration-300 ${c.cover_sensitive ? 'scale-125 blur-2xl' : 'group-hover:scale-105'}`}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-brand-100 to-brand-200 text-7xl font-extrabold text-brand-600rand-900rand-600">
            {c.title.charAt(0).toUpperCase()}
          </div>
        )}
        {c.cover_url && c.cover_sensitive && (
          <span className="absolute left-1/2 top-1/2 inline-flex -translate-x-1/2 -translate-y-1/2 items-center gap-2 rounded-full bg-[#111827]/80 px-4 py-2 text-sm font-semibold text-white">
            <EyeOff size={16} /> Image sensible
          </span>
        )}
        {country && (
          <span className="absolute left-3 top-3 rounded-full bg-white px-3 py-1 text-xs font-bold text-brand-600rand-900rand-700">{country}</span>
        )}
      </div>

      <div className="p-5">
        <h3 className={`line-clamp-2 font-bold leading-snug text-[#111827] ${big ? 'text-xl' : 'text-lg'}`}>{c.title}</h3>
        {goal > 0 && (
          <div className="mt-4 h-2 overflow-hidden rounded-full bg-gray-200">
            <div className="h-full rounded-full bg-brand-600rand-500" style={{ width: `${pct}%` }} />
          </div>
        )}
        <div className={`flex items-baseline justify-between text-sm text-gray-600 ${goal > 0 ? 'mt-3' : 'mt-4'}`}>
          <span><b className="text-[#111827]">{fcfa(c.current_amount)}</b> récoltés</span>
          {goal > 0 && <span>{pct >= 100 ? 'Objectif atteint' : `${pct} %`}</span>}
        </div>
      </div>
    </Link>
  )
}