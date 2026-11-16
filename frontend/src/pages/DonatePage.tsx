import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, fcfa } from '../api'
import { PAY_METHODS, calcFee } from '../constants'
import type { CampaignDetail, DonationCreated } from '../types'

const QUICK = [2000, 5000, 10000, 25000]
const PHONE_RE = /^\+\d{8,15}$/

export default function DonatePage() {

  const { slug } = useParams<{ slug: string }>()
  const [campaign, setCampaign] = useState<CampaignDetail | null>(null)
  const [amount, setAmount] = useState('')
  const [method, setMethod] = useState('')
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')
  const [email, setEmail] = useState('')
  const [comment, setComment] = useState('')
  const [anonymous, setAnonymous] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [result, setResult] = useState<DonationCreated | null>(null)

  useEffect(() => {
    api<CampaignDetail>(`/campaigns/${slug}`).then(setCampaign).catch((e: Error) => setError(e.message))
  }, [slug])

  const amountNum = parseInt(amount.replace(/\D/g, ''), 10) || 0
  const fee = calcFee(amountNum)
  const total = amountNum + fee
  const selected = PAY_METHODS.find((m) => m.value === method)
  const overLimit = !!selected && selected.max > 0 && total > selected.max
  const cleanPhone = phone.replace(/[\s.-]/g, '')

  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    if (!campaign) return
    if (amountNum <= 0) return setError('Choisis un montant.')
    if (!method) return setError('Choisis un moyen de paiement.')
    if (overLimit) return setError('Ce montant dépasse la limite de ce moyen de paiement.')
    if (cleanPhone && !PHONE_RE.test(cleanPhone)) {
      return setError('Numéro au format international, par exemple +223 70 00 00 00.')
    }
    setBusy(true)
    try {
      const d = await api<DonationCreated>('/donations', {
        method: 'POST',
        body: JSON.stringify({
          campaign_id: campaign.id,
          amount: amountNum,
          tip_amount: 0,
          payment_method: method,
          donor_name: name,
          donor_phone: cleanPhone,
          donor_email: email,
          comment,
          is_anonymous: anonymous,
        }),
      })
      setResult(d)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setBusy(false)
    }
  }

  if (!campaign) return <p className="py-10">{error || 'Chargement…'}</p>

  if (result) {
    return (
      <div className="mx-auto max-w-lg py-12 text-brand-600enter">
        <p className="text-brand-600rand-900xl font-bold text-[#111827]">Don enregistré</p>
        <p className="mt-3 text-gray-600">
          Total à payer : <span className="font-bold">{fcfa(result.total_charged)}</span>
        </p>
        <p className="mt-4 rounded-2xl bg-brand-600mber-50 p-4 text-sm text-brand-600rand-900mber-800">
          Le paiement en ligne n'est pas encore activé. Ce don est en attente.
        </p>
        <Link to={`/c/${campaign.url_slug}`} className="mt-6 inline-block rounded-full bg-[#111827] px-8 py-4 text-lg font-semibold text-white">
          Retour à la cagnotte
        </Link>
      </div>
    )
  }

  return (
    <form onSubmit={submit} className="mx-auto max-w-lg space-y-6 py-8">
      <div>
        <p className="text-sm text-gray-500">Je soutiens</p>
        <h1 className="text-brand-600rand-900xl font-bold text-[#111827]">{campaign.title}</h1>
      </div>

      <section>
        <p className="mb-2 text-lg font-semibold">Montant</p>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
          {QUICK.map((q) => (
            <button
              type="button"
              key={q}
              onClick={() => setAmount(String(q))}
              className={`rounded-2xl border-2 py-3 font-semibold ${amountNum === q ? 'border-brand-600 bg-brand-600rand-50 text-brand-600rand-900rand-700' : 'border-gray-200 bg-white'}`}
            >
              {new Intl.NumberFormat('fr-FR').format(q)}
            </button>
          ))}
        </div>
        <div className="relative mt-3">
          <input
            inputMode="numeric"
            value={amount}
            onChange={(e) => setAmount(e.target.value.replace(/\D/g, ''))}
            placeholder="Autre montant"
            className="w-full rounded-2xl border border-gray-300 bg-white p-4 pr-20 text-xl outline-none focus:border-brand-500"
          />
          <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500">FCFA</span>
        </div>
      </section>

      <section>
        <p className="mb-2 text-lg font-semibold">Payer avec</p>
        <div className="grid grid-cols-2 gap-3">
          {PAY_METHODS.map((m) => (
            <button
              type="button"
              key={m.value}
              onClick={() => setMethod(m.value)}
              className={`flex h-20 flex-col items-center justify-center gap-1 rounded-2xl border-2 bg-white ${method === m.value ? 'border-brand-600 ring-2 ring-brand-200' : 'border-gray-200'}`}
            >
              <img src={m.logo} alt="" className="h-8 w-auto object-contain" onError={(e) => (e.currentTarget.style.display = 'none')} />
              <span className="text-sm font-medium">{m.label}</span>
            </button>
          ))}
        </div>
        {overLimit && selected && (
          <p className="mt-2 text-red-600">
            {selected.label} accepte au maximum {fcfa(selected.max)} par paiement, frais compris.
          </p>
        )}
      </section>

      <section className="space-y-3">
        <p className="text-lg font-semibold">Tes informations <span className="text-sm font-normal text-gray-500">(facultatif)</span></p>
        <input value={name} onChange={(e) => setName(e.target.value)} maxLength={100} placeholder="Ton prénom" className="w-full rounded-2xl border border-gray-300 bg-white p-4 text-lg outline-none focus:border-brand-500" />
        <input value={phone} onChange={(e) => setPhone(e.target.value)} inputMode="tel" placeholder="Téléphone, ex. +223 70 00 00 00" className="w-full rounded-2xl border border-gray-300 bg-white p-4 text-lg outline-none focus:border-brand-500" />
        <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="Email" className="w-full rounded-2xl border border-gray-300 bg-white p-4 text-lg outline-none focus:border-brand-500" />
        <textarea value={comment} onChange={(e) => setComment(e.target.value)} maxLength={500} rows={3} placeholder="Un message d'encouragement (visible sur la cagnotte)" className="w-full rounded-2xl border border-gray-300 bg-white p-4 text-lg outline-none focus:border-brand-500" />
        <label className="flex items-center gap-3 text-brand-600rand-900ase">
          <input type="checkbox" checked={anonymous} onChange={(e) => setAnonymous(e.target.checked)} className="h-5 w-5" />
          Faire ce don en anonyme
        </label>
      </section>

      <section className="rounded-2xl bg-gray-50 p-4">
        <div className="flex justify-between text-gray-600"><span>Ton don</span><span>{fcfa(amountNum)}</span></div>
        <div className="mt-1 flex justify-between text-sm text-gray-500"><span>Frais Demso</span><span>{fcfa(fee)}</span></div>
        <div className="mt-2 flex justify-between border-t border-gray-200 pt-2 text-lg font-bold"><span>Total</span><span>{fcfa(total)}</span></div>
      </section>

      {error && <p className="text-red-600">{error}</p>}

      <button disabled={busy} className="w-full rounded-full bg-brand-600rand-500 py-4 text-xl font-semibold text-[#111827] disabled:opacity-60">
        {busy ? 'Envoi…' : 'Continuer'}
      </button>
    </form>
  )
}