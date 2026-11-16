import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, fcfa } from '../api'
import type { Campaign } from '../types'

export default function MyCampaignsPage() {
  const [list, setList] = useState<Campaign[] | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    api<{ campaigns: Campaign[] }>('/me/campaigns')
      .then((r) => setList(r.campaigns))
      .catch((e: Error) => setError(e.message))
  }, [])

  return (
    <div className="py-6">
      <div className="flex items-center justify-between gap-3">
        <h1 className="text-3xl font-extrabold text-[#111827]">Mes cagnottes</h1>
        <Link to="/creer" className="inline-flex items-center rounded-full border-2 border-gray-900 px-5 py-2.5 text-sm font-semibold text-gray-900 transition hover:bg-gray-900 hover:text-white">
          Nouvelle cagnotte
        </Link>
      </div>

      {error && <p className="mt-6 text-red-600">{error}</p>}
      {!list && !error && <p className="mt-6 text-gray-500">Chargement…</p>}

      {list && list.length === 0 && (
        <div className="rounded-3xl border-2 border-dashed border-brand-200 bg-brand-50 p-8 text-center">
          <p className="text-lg text-gray-700">Tu n'as pas encore de cagnotte.</p>
          <Link to="/creer"className="btn btn-primary mt-6">
            Créer ma première cagnotte
          </Link>
        </div>
      )}

      <div className="mt-6 space-y-4">
        {list?.map((c) => {
          const pct = c.goal_amount ? Math.min(100, Math.round((c.current_amount / c.goal_amount) * 100)) : 0
          return (
            <div key={c.id} className="rounded-2xl border border-gray-200 bg-white p-5">
              <p className="text-lg font-bold text-[#111827]">{c.title}</p>
              <p className="mt-1 text-gray-600">
                {fcfa(c.current_amount)} récoltés{c.goal_amount ? ` sur ${fcfa(c.goal_amount)}` : ''}
              </p>
              {c.goal_amount ? (
                <div className="mt-3 h-2 overflow-hidden rounded-full bg-gray-200">
                  <div className="h-full rounded-full bg-brand-600rand-500" style={{ width: `${pct}%` }} />
                </div>
              ) : null}
              <div className="mt-4 flex gap-5">
                <Link to={`/c/${c.url_slug}`} className="font-semibold text-brand-600rand-900rand-700">Voir la cagnotte →</Link>
                <Link to={`/mon-espace/c/${c.id}/messages`} className="font-semibold text-gray-700">Messages</Link>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}