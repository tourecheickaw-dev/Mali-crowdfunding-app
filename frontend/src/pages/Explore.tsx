import { useEffect, useMemo, useState } from 'react'
import { Search } from 'lucide-react'
import { api } from '../api'
import CampaignCard from '../components/CampaignCard'
import type { Campaign, CampaignList } from '../types'
import { Link } from 'react-router-dom'

export default function Explore() {
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [query, setQuery] = useState('')

  useEffect(() => {
    api<CampaignList | Campaign[]>('/campaigns')
      .then((d) => setCampaigns(Array.isArray(d) ? d : d.campaigns ?? []))
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  const shown = useMemo(() => {
    const q = query.trim().toLowerCase()
    return q ? campaigns.filter((c) => c.title.toLowerCase().includes(q)) : campaigns
  }, [campaigns, query])

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold text-[#111827]">Toutes les cagnottes</h1>

      <div className="relative mt-5 max-w-md">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Rechercher une cagnotte"
          className="w-full rounded-full border border-gray-300 bg-white py-3 pl-12 pr-4 text-lg outline-none focus:border-brand-500"
        />
      </div>

      <div className="mt-6">
        {loading && <p>Chargement…</p>}
        {error && <p className="text-red-600">{error}</p>}
        {!loading && !error && shown.length === 0 && (
  <div className="rounded-3xl border-2 border-dashed border-brand-200 bg-brand-50 p-8 text-center">
    <p className="text-xl font-bold text-gray-900">La première cagnotte commence ici.</p>
    <p className="mt-2 text-gray-600">Ça prend quelques minutes.</p>
    <Link to="/creer" className="btn btn-primary mt-5">
      Créer ma cagnotte
    </Link>
  </div>
)}
        
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {shown.map((c) => <CampaignCard key={c.id} c={c} />)}
        </div>
      </div>
    </div>
  )
}