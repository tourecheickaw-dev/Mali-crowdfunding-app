import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api'
import CampaignCard from '../components/CampaignCard'
import HowItWorks from '../components/HowItWorks'
import PaymentLogos from '../components/PaymentLogos'
import TrustBar from '../components/TrustBar'
import type { Campaign, CampaignList } from '../types'

export default function Home() {
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    api<CampaignList | Campaign[]>('/campaigns')
      .then((d) => setCampaigns(Array.isArray(d) ? d : d.campaigns ?? []))
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <section className="pt-16 text-brand-600enter">
     <h1 className="text-center text-4xl font-extrabold leading-[1.1] tracking-tight text-brand-800 sm:text-6xl">
  Ensemble, <span className="text-brand-600">on va plus loin.</span>
</h1>
<p className="mx-auto mt-5 max-w-2xl text-center text-lg text-gray-600">
  Aidez ceux que vous aimez, où que vous soyez. Payez avec Wave, Orange Money ou carte.
</p>

       <div className="mt-12 flex flex-col items-center justify-center gap-4 sm:flex-row sm:gap-6">
  <Link to="/cagnottes" className="btn btn-outline">Je veux aider</Link>
  <Link to="/creer" className="btn btn-primary">Créer ma cagnotte</Link>
</div>

        <p className="mx-auto mt-8 max-w-xl text-center text-sm text-gray-500">
          Une urgence, un projet, un proche à soutenir ? Crée ta cagnotte en quelques minutes et partage-la à ta
          famille et tes amis.
        </p>
      </section>

      <section id="cagnottes" className="mt-20">
        <h2 className="text-2xl font-extrabold tracking-tight text-gray-900">Cagnottes à soutenir</h2>

        {loading && <p className="mt-4 text-gray-500">Chargement…</p>}
        {error && <p className="mt-4 text-red-600">{error}</p>}

       {!loading && !error && campaigns.length === 0 && (
  <div className="mt-6 rounded-3xl border-2 border-dashed border-brand-200 bg-brand-50 p-8 text-center">
    <p className="text-xl font-extrabold text-gray-900">La première cagnotte commence ici.</p>
    <p className="mt-2 text-gray-600">Ça prend quelques minutes.</p>
    <Link to="/creer" className="btn btn-primary mt-6">
      Créer ma cagnotte
    </Link>
  </div>
)}

        {campaigns.length > 0 && (
          <div className="mt-6 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {campaigns.slice(0, 6).map((c) => (
              <CampaignCard key={c.id} c={c} />
            ))}
          </div>
        )}

        <div className="mt-10 text-center">
  <Link to="/cagnottes" className="btn btn-outline">
    Voir toutes les cagnottes
  </Link>
</div>
      </section>

      <div className="mt-20">
        <HowItWorks />
      </div>
      <div className="mt-16">
        <PaymentLogos />
      </div>
      <TrustBar />
    </div>
  )
}