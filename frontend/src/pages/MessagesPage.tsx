import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, fcfa } from '../api'
import { formatDate } from '../format'
import type { Campaign, OrganizerDonation } from '../types'

function Row({ campaignId, d, onChange }: { campaignId: string; d: OrganizerDonation; onChange: (n: OrganizerDonation) => void }) {
  const [text, setText] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const base = `/me/campaigns/${campaignId}/donations/${d.id}`

  async function run(action: () => Promise<void>) {
    setBusy(true)
    setError('')
    try {
      await action()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Erreur')
    } finally {
      setBusy(false)
    }
  }

  const sendReply = () =>
    run(async () => {
      const reply = text.trim()
      if (!reply) throw new Error('Écris ta réponse.')
      await api(`${base}/reply`, { method: 'PUT', body: JSON.stringify({ reply }) })
      onChange({ ...d, reply, reply_at: new Date().toISOString() })
      setText('')
    })

  const removeReply = () =>
    run(async () => {
      await api(`${base}/reply`, { method: 'DELETE' })
      onChange({ ...d, reply: undefined, reply_at: undefined })
    })

  const toggleHidden = () =>
    run(async () => {
      await api(`${base}/comment-visibility`, { method: 'PUT', body: JSON.stringify({ hidden: !d.comment_hidden }) })
      onChange({ ...d, comment_hidden: !d.comment_hidden })
    })

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5">
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="font-bold text-[#111827]">
            {d.donor_name || 'Un donateur'}
            {d.is_anonymous && <span className="ml-2 rounded-full bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-600">anonyme pour le public</span>}
          </p>
          <p className="text-sm text-gray-500">{formatDate(d.created_at)}</p>
        </div>
        <p className="font-bold text-brand-600rand-900rand-700">{fcfa(d.amount)}</p>
      </div>

      {d.comment && (
        <div className="mt-3">
          <p className={`whitespace-pre-wrap break-words text-gray-800 ${d.comment_hidden ? 'opacity-50' : ''}`}>{d.comment}</p>
          {d.comment_hidden && <p className="mt-1 text-sm font-medium text-brand-600rand-900mber-700">Ce message est masqué pour le public.</p>}
          <button type="button" onClick={toggleHidden} disabled={busy} className="mt-2 text-sm font-semibold text-gray-600 underline">
            {d.comment_hidden ? 'Afficher ce message' : 'Masquer ce message'}
          </button>
        </div>
      )}

      {d.comment && d.reply && (
        <div className="mt-3 rounded-2xl bg-brand-600rand-50 p-3">
          <p className="text-sm font-semibold text-brand-600rand-900rand-800">Ta réponse</p>
          <p className="whitespace-pre-wrap break-words text-gray-800">{d.reply}</p>
          <button type="button" onClick={removeReply} disabled={busy} className="mt-1 text-sm font-semibold text-red-600">
            Supprimer ma réponse
          </button>
        </div>
      )}

      {d.comment && !d.reply && (
        <div className="mt-3 space-y-2">
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            maxLength={500}
            rows={2}
            placeholder="Remercie cette personne…"
            className="w-full rounded-2xl border-2 border-gray-200 bg-white p-3 text-[#111827] outline-none focus:border-brand-500"
          />
          <button type="button" onClick={sendReply} disabled={busy} className="rounded-full bg-brand-600rand-500 px-5 py-2 font-semibold text-[#111827] disabled:opacity-60">
            Répondre
          </button>
        </div>
      )}

      {error && <p className="mt-2 text-red-600">{error}</p>}
    </div>
  )
}

export default function MessagesPage() {
  const { id } = useParams<{ id: string }>()
  const [donations, setDonations] = useState<OrganizerDonation[] | null>(null)
  const [title, setTitle] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([
      api<OrganizerDonation[]>(`/me/campaigns/${id}/donations`),
      api<{ campaigns: Campaign[] }>('/me/campaigns'),
    ])
      .then(([list, mine]) => {
        setDonations(list)
        setTitle(mine.campaigns.find((c) => c.id === id)?.title ?? '')
      })
      .catch((e: Error) => setError(e.message))
  }, [id])

  function update(n: OrganizerDonation) {
    setDonations((prev) => (prev ? prev.map((x) => (x.id === n.id ? n : x)) : prev))
  }

  return (
    <div className="py-6">
      <Link to="/mon-espace" className="text-sm font-semibold text-gray-600">← Mes cagnottes</Link>
      <h1 className="mt-2 text-brand-600rand-900xl font-extrabold text-[#111827]">Messages et dons</h1>
      {title && <p className="text-gray-600">{title}</p>}

      {error && <p className="mt-6 text-red-600">{error}</p>}
      {!donations && !error && <p className="mt-6 text-gray-500">Chargement…</p>}

      {donations && donations.length === 0 && (
        <p className="mt-6 rounded-2xl bg-gray-50 p-6 text-gray-600">
          Aucun don confirmé pour l'instant. Les dons apparaissent ici dès que le paiement est confirmé.
        </p>
      )}

      <div className="mt-6 space-y-4">
        {donations?.map((d) => (
          <Row key={d.id} campaignId={id ?? ''} d={d} onChange={update} />
        ))}
      </div>
    </div>
  )
}