import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, fcfa } from '../api'
import PaymentLogos from '../components/PaymentLogos'
import type { CampaignDetail, PublicDonation } from '../types'
import ShareButtons from '../components/ShareButtons'

type Tab = 'description' | 'messages' | 'contributions'

export default function CampaignPage() {
  const { slug } = useParams<{ slug: string }>()
  const [campaign, setCampaign] = useState<CampaignDetail | null>(null)
  const [donations, setDonations] = useState<PublicDonation[]>([])
  const [tab, setTab] = useState<Tab>('description')
  const [error, setError] = useState('')

  useEffect(() => {
    api<CampaignDetail>(`/campaigns/${slug}`)
      .then(setCampaign)
      .catch((e: Error) => setError(e.message))

    api<{ donations: PublicDonation[] } | PublicDonation[]>(`/campaigns/${slug}/donations`)
      .then((d) => setDonations(Array.isArray(d) ? d : d.donations ?? []))
      .catch(() => setDonations([]))
  }, [slug])

  if (error) return <p className="py-10 text-red-600">{error}</p>
  if (!campaign) return <p className="py-10 text-gray-700">Chargement…</p>

  const goal = campaign.goal_amount
  const pct = goal ? Math.min(100, Math.round((campaign.current_amount / goal) * 100)) : null
  const messages = donations.filter((d) => d.comment)

  const tabs: { id: Tab; label: string }[] = [
    { id: 'description', label: 'Description' },
    { id: 'messages', label: `Messages (${messages.length})` },
    { id: 'contributions', label: `Dons (${donations.length})` },
  ]

  return (
    <div className="py-6">
      <div className="grid gap-8 md:grid-cols-3">
        <div className="md:col-span-2">
          <div className="h-64 overflow-hidden rounded-3xl bg-gray-200 sm:h-80">
            {campaign.cover_url && <img src={campaign.cover_url} alt="" className="h-full w-full object-cover" />}
          </div>
          <h1 className="mt-5 text-2xl font-extrabold leading-tight tracking-tight text-gray-900 sm:text-4xl">
            {campaign.title}
          </h1>
          {campaign.organizer_name && (
            <p className="mt-2 text-gray-600">
              Organisée par <span className="font-semibold text-gray-900">{campaign.organizer_name}</span>
            </p>
          )}
        </div>

        <div className="h-fit rounded-3xl bg-white p-6 shadow-[0_10px_40px_-15px_rgba(0,0,0,0.2)] ring-1 ring-gray-200 md:sticky md:top-24">
          <p className="text-3xl font-extrabold tracking-tight text-gray-900">{fcfa(campaign.current_amount)}</p>
          {goal && <p className="mt-1 text-gray-600">sur {fcfa(goal)}</p>}

          {pct !== null && (
            <div className="mt-4 h-3 overflow-hidden rounded-full bg-gray-200">
              <div
                className="h-3 rounded-full transition-all"
                style={{ width: pct + '%', backgroundColor: '#0e7c4a' }}
              />
            </div>
          )}
          {pct === 100 && (
            <p className="mt-3 inline-block rounded-full bg-[#fdf1cc] px-3 py-1 text-sm font-bold text-gray-900">
              🎉 Objectif atteint
            </p>
          )}
          <p className="mt-3 text-sm text-gray-600">{donations.length} don(s)</p>

          <Link to={`/c/${campaign.url_slug}/don`} className="btn btn-primary mt-5 w-full">
            Faire un don
          </Link>

          <div className="mt-6 border-t border-gray-200 pt-5">
            <p className="mb-3 font-extrabold text-gray-900">Partager cette cagnotte</p>
            <ShareButtons title={campaign.title} />
          </div>

          <div className="mt-6">
            <PaymentLogos />
          </div>
        </div>
      </div>

      <div className="mt-10 flex gap-2 border-b border-gray-200">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => setTab(t.id)}
            className={`-mb-px px-4 py-3 text-base font-semibold transition ${
              tab === t.id ? 'border-b-4 border-brand-600 text-gray-900' : 'border-b-4 border-transparent text-gray-500 hover:text-gray-900'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className="mt-6">
        {tab === 'description' && (
          <p className="whitespace-pre-line text-lg leading-relaxed text-gray-900">
            {campaign.description || 'Pas de description.'}
          </p>
        )}

        {tab === 'messages' && (
          <div className="space-y-4">
            {messages.length === 0 && <p className="text-gray-600">Aucun message pour le moment.</p>}
            {messages.map((d) => (
              <div key={d.id} className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-gray-200">
                <p className="font-bold text-gray-900">{d.donor_name || 'Anonyme'}</p>
                <p className="mt-1 whitespace-pre-line text-gray-800">{d.comment}</p>
                {d.reply && (
                  <div className="ml-2 mt-4 rounded-r-2xl border-l-4 border-brand-600 bg-brand-50 p-4">
                    <p className="text-sm font-bold text-gray-900">Réponse de l'organisateur</p>
                    <p className="mt-1 whitespace-pre-line text-gray-800">{d.reply}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}

        {tab === 'contributions' && (
          <div className="space-y-3">
            {donations.length === 0 && <p className="text-gray-600">Soyez le premier à donner.</p>}
            {donations.map((d) => (
              <div
                key={d.id}
                className="flex items-center justify-between rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-200"
              >
                <span className="font-semibold text-gray-900">{d.donor_name || 'Anonyme'}</span>
                <span className="font-extrabold text-gray-900">{fcfa(d.amount)}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
