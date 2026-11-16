import { useCallback, useEffect, useState } from 'react'
import { api } from '../api'
import { PAYOUT_TYPES, UEMOA_COUNTRIES } from '../constants'
import { maskNumber } from '../format'
import type { PayoutAccount } from '../types'

const PHONE_RE = /^\+\d{8,15}$/
const input = 'w-full rounded-2xl border-2 border-gray-200 bg-white p-4 text-lg text-[#111827] outline-none focus:border-brand-500'
const labelOf = (t: string): string => PAYOUT_TYPES.find((p) => p.value === t)?.label ?? t

export default function PayoutPage() {
  const [accounts, setAccounts] = useState<PayoutAccount[] | null>(null)
  const [error, setError] = useState('')
  const [busyId, setBusyId] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [type, setType] = useState('')
  const [name, setName] = useState('')
  const [number, setNumber] = useState('')
  const [bankName, setBankName] = useState('')
  const [country, setCountry] = useState('ML')
  const [saving, setSaving] = useState(false)

  const load = useCallback(() => {
    api<PayoutAccount[]>('/me/payout-accounts')
      .then(setAccounts)
      .catch((e: Error) => setError(e.message))
  }, [])

  useEffect(() => {
    load()
  }, [load])

  async function activate(a: PayoutAccount) {
    const ok = window.confirm(`Les prochains dons iront sur ${labelOf(a.type)} ${maskNumber(a.account_number)}. Continuer ?`)
    if (!ok) return
    setBusyId(a.id)
    setError('')
    try {
      await api(`/me/payout-accounts/${a.id}/activate`, { method: 'POST' })
      load()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Erreur')
    } finally {
      setBusyId('')
    }
  }

  async function add(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    const isBank = type === 'bank'
    const clean = number.replace(/[\s.-]/g, '')
    if (!type) return setError('Choisis un type de compte.')
    if (!name.trim()) return setError('Écris le nom du titulaire.')
    if (isBank && !bankName.trim()) return setError('Écris le nom de la banque.')
    if (!isBank && !PHONE_RE.test(clean)) return setError('Numéro au format international, par exemple +223 70 00 00 00.')
    if (isBank && clean.length < 8) return setError('Numéro de compte trop court.')
    setSaving(true)
    try {
      await api('/me/payout-accounts', {
        method: 'POST',
        body: JSON.stringify({
          type,
          account_name: name.trim(),
          account_number: clean,
          bank_name: isBank ? bankName.trim() : '',
          country,
        }),
      })
      setType('')
      setName('')
      setNumber('')
      setBankName('')
      setShowForm(false)
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="py-6">
      <h1 className="text-brand-600rand-900xl font-extrabold text-[#111827]">Mes comptes de paiement</h1>
      <p className="mt-1 text-gray-600">Un seul compte reçoit les dons à la fois : le compte actif.</p>

      {import.meta.env.DEV && (
        <p className="mt-4 rounded-2xl bg-brand-600mber-50 p-3 text-sm text-brand-600rand-900mber-800">
          Mode test : la vérification de sécurité pour changer de compte (code, délai de 24 à 48 h) n'est pas encore active.
        </p>
      )}

      {error && <p className="mt-4 text-red-600">{error}</p>}
      {!accounts && !error && <p className="mt-6 text-gray-500">Chargement…</p>}

      <div className="mt-6 space-y-3">
        {accounts?.map((a) => (
          <div key={a.id} className={`rounded-2xl border-2 bg-white p-5 ${a.is_active ? 'border-brand-500' : 'border-gray-200'}`}>
            <div className="flex items-center gap-3">
              {PAYOUT_TYPES.find((p) => p.value === a.type)?.logo && (
                <img src={PAYOUT_TYPES.find((p) => p.value === a.type)?.logo} alt="" className="h-8 w-auto object-contain" onError={(e) => (e.currentTarget.style.display = 'none')} />
              )}
              <p className="text-lg font-bold text-[#111827]">{labelOf(a.type)}</p>
              {a.is_active && <span className="rounded-full bg-brand-600rand-100 px-3 py-1 text-sm font-semibold text-brand-600rand-900rand-800">Actif</span>}
            </div>
            <p className="mt-2 text-gray-800">{a.account_name}</p>
            <p className="text-gray-600">
              {maskNumber(a.account_number)}
              {a.bank_name ? ` · ${a.bank_name}` : ''} · {a.country}
            </p>
            <p className="mt-1 text-sm text-gray-500">{a.verified_at ? 'Vérifié' : 'Non vérifié'}</p>
            {a.is_active ? (
              <p className="mt-3 text-sm font-medium text-brand-600rand-900rand-700">Les prochains dons arrivent sur ce compte.</p>
            ) : (
              <button type="button" disabled={busyId === a.id} onClick={() => activate(a)} className="mt-3 rounded-full bg-[#111827] px-5 py-2 font-semibold text-white disabled:opacity-60">
                {busyId === a.id ? 'Changement…' : 'Utiliser pour les prochains dons'}
              </button>
            )}
          </div>
        ))}
      </div>

      {accounts && accounts.length === 0 && <p className="mt-6 rounded-2xl bg-gray-50 p-6 text-gray-600">Aucun compte pour l'instant.</p>}

      {!showForm ? (
        <button type="button" onClick={() => setShowForm(true)} className="mt-6 rounded-full bg-brand-600rand-500 px-6 py-3 font-semibold text-[#111827]">
          Ajouter un compte
        </button>
      ) : (
        <form onSubmit={add} className="mt-6 space-y-3 rounded-2xl border border-gray-200 bg-white p-5">
          <div className="grid grid-cols-2 gap-3">
            {PAYOUT_TYPES.map((p) => (
              <button type="button" key={p.value} onClick={() => setType(p.value)}
                className={`rounded-2xl border-2 p-3 font-medium ${type === p.value ? 'border-brand-600 bg-brand-600rand-50' : 'border-gray-200'}`}>
                {p.label}
              </button>
            ))}
          </div>
          {type && (
            <>
              <input value={name} onChange={(e) => setName(e.target.value)} maxLength={100} placeholder="Nom du titulaire" className={input} />
              {type === 'bank' && <input value={bankName} onChange={(e) => setBankName(e.target.value)} maxLength={100} placeholder="Nom de la banque" className={input} />}
              <input value={number} onChange={(e) => setNumber(e.target.value)} placeholder={type === 'bank' ? 'Numéro de compte (IBAN ou RIB)' : 'Numéro, ex. +223 70 00 00 00'} className={input} />
              <select value={country} onChange={(e) => setCountry(e.target.value)} className={input}>
                {UEMOA_COUNTRIES.map((c) => <option key={c.code} value={c.code}>{c.name}</option>)}
              </select>
            </>
          )}
          <div className="flex gap-3">
            <button disabled={saving} className="rounded-full bg-brand-600rand-500 px-6 py-3 font-semibold text-[#111827] disabled:opacity-60">
              {saving ? 'Ajout…' : 'Ajouter'}
            </button>
            <button type="button" onClick={() => setShowForm(false)} className="rounded-full px-6 py-3 font-semibold text-gray-700">Annuler</button>
          </div>
        </form>
      )}
    </div>
  )
}