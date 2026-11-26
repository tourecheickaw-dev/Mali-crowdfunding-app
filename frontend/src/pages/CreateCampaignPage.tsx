import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api'
import { BENEFICIARY_TYPES, PAYOUT_TYPES, UEMOA_COUNTRIES } from '../constants'
import PhotoPicker from '../components/PhotoPicker'
import { attachCover } from '../upload'
import type { Campaign } from '../types'
import { Landmark, ArrowRight, Rocket } from 'lucide-react'

type PickedPhoto = { url: string; sensitive: boolean }

const PHONE_RE = /^\+\d{8,15}$/
const input =
  'w-full rounded-2xl border-2 border-gray-200 bg-white p-4 text-lg text-gray-900 outline-none focus:border-brand-600'
const smallInput =
  'w-full rounded-2xl border-2 border-gray-200 bg-white px-4 py-3 text-base text-gray-900 outline-none transition focus:border-brand-600'

const HEADINGS = [
  'Pour qui collectes-tu ?',
  'Donne un nom à ta cagnotte',
  'Raconte ton histoire',
  "Où veux-tu recevoir l'argent ?",
]

const HINTS = [
  'Chaque grande aide commence par un premier pas.',
  'Un titre clair donne envie de lire la suite.',
  'Les gens donnent quand ils comprennent. Sois simple et sincère.',
  "Dernière étape : tout l'argent des dons arrivera directement ici.",
]

// Messages en français quand l'ajout de la banque est refusé par le serveur
function bankErrorMessage(msg: string): string {
  if (msg.includes('verified mobile money')) {
    return "Ton compte mobile money doit d'abord être vérifié. Tu pourras ajouter ta banque ensuite, depuis « Mes comptes de paiement »."
  }
  if (msg.includes('must match')) {
    return 'Le nom du compte bancaire doit être le même que celui de ton compte mobile money. Tu pourras réessayer depuis « Mes comptes de paiement ».'
  }
  return "Le compte bancaire n'a pas pu être ajouté. Tu pourras réessayer depuis « Mes comptes de paiement »."
}

export default function CreateCampaignPage() {
  const nav = useNavigate()
  const [step, setStep] = useState(0)
  const [hasPayout, setHasPayout] = useState(false)

  const [benType, setBenType] = useState('self')
  const [benName, setBenName] = useState('')
  const [country, setCountry] = useState('ML')
  const [title, setTitle] = useState('')
  const [goal, setGoal] = useState('')
  const [description, setDescription] = useState('')
  const [photo, setPhoto] = useState<PickedPhoto | null>(null)

  // main payout account (mobile money only)
  const [payType, setPayType] = useState('')
  const [payName, setPayName] = useState('')
  const [payNumber, setPayNumber] = useState('')
  const [bankTapped, setBankTapped] = useState(false)

  // optional bank account, added after the mobile money account
  const [addBank, setAddBank] = useState(false)
  const [extraBankName, setExtraBankName] = useState('')
  const [extraHolder, setExtraHolder] = useState('')
  const [extraNumber, setExtraNumber] = useState('')

  const [payoutCreated, setPayoutCreated] = useState(false)
  const [bankCreated, setBankCreated] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [done, setDone] = useState<{ slug: string; notes: string[] } | null>(null)

  const goalNum = parseInt(goal.replace(/\D/g, ''), 10) || 0
  const bankRecommended = goalNum >= 1_000_000 || goalNum === 0

  // Skip the payout step if the organizer already has an account.
  useEffect(() => {
    api<unknown>('/me/payout-accounts')
      .then((res) => {
        const list = Array.isArray(res)
          ? res
          : Object.values((res ?? {}) as Record<string, unknown>).find((v) => Array.isArray(v))
        setHasPayout(Array.isArray(list) && list.length > 0)
      })
      .catch(() => setHasPayout(false))
  }, [])

  const total = hasPayout ? 3 : 4
  const cleanNumber = payNumber.replace(/[\s.-]/g, '')
  const cleanExtra = extraNumber.replace(/\s/g, '')

  function pickPayType(value: string) {
    if (value === 'bank') {
      // La banque vient toujours après un mobile money
      if (payType) setAddBank(true)
      else setBankTapped(true)
      return
    }
    setBankTapped(false)
    setPayType(value)
  }

  function validate(s: number): string {
    if (s === 0 && benType !== 'self' && !benName.trim()) return 'Écris le nom de la personne ou de la cause.'
    if (s === 1 && !title.trim()) return 'Donne un titre à ta cagnotte.'
    if (s === 2 && description.trim().length < 20)
      return 'Raconte ton histoire en quelques phrases (20 caractères minimum).'
    if (s === 3) {
      if (!payType) return "Choisis où recevoir l'argent."
      if (!payName.trim()) return 'Écris le nom du titulaire du compte.'
      if (!PHONE_RE.test(cleanNumber)) return 'Numéro au format international, par exemple +223 70 00 00 00.'
      if (addBank) {
        if (!extraBankName.trim()) return 'Écris le nom de la banque.'
        if (!extraHolder.trim()) return 'Écris le nom du titulaire du compte bancaire.'
        if (cleanExtra.length < 8) return 'Numéro du compte bancaire trop court.'
      }
    }
    return ''
  }

  function next() {
    const msg = validate(step)
    setError(msg)
    if (msg) return
    if (step < total - 1) setStep(step + 1)
    else void submit()
  }

  async function submit() {
    setBusy(true)
    setError('')
    try {
      if (!hasPayout && !payoutCreated) {
        const acc = await api<{ id?: string }>('/me/payout-accounts', {
          method: 'POST',
          body: JSON.stringify({
            type: payType,
            account_name: payName.trim(),
            account_number: cleanNumber,
            bank_name: '',
            country,
          }),
        })
        setPayoutCreated(true)
        if (acc.id) {
          try {
            await api(`/me/payout-accounts/${acc.id}/activate`, { method: 'POST' })
          } catch {
            // the account may already be active; the campaign call will tell us if not
          }
        }
      }

      const c = await api<Campaign>('/me/campaigns', {
        method: 'POST',
        body: JSON.stringify({
          title: title.trim(),
          description,
          goal_amount: goalNum > 0 ? goalNum : null,
          country,
          beneficiary_type: benType,
          beneficiary_name: benType === 'self' ? '' : benName.trim(),
        }),
      })

      const notes: string[] = []

      if (photo) {
        try {
          await attachCover(c.id, photo.url, photo.sensitive)
        } catch {
          notes.push("La photo n'a pas pu être ajoutée. Tu pourras réessayer plus tard.")
        }
      }

      // Bank account: created after the campaign, never blocks it.
      // It is NOT activated: donations keep going to the main account.
      if (!hasPayout && addBank && !bankCreated) {
        try {
          await api('/me/payout-accounts', {
            method: 'POST',
            body: JSON.stringify({
              type: 'bank',
              account_name: extraHolder.trim(),
              account_number: cleanExtra,
              bank_name: extraBankName.trim(),
              country,
            }),
          })
          setBankCreated(true)
        } catch (e) {
          notes.push(bankErrorMessage(e instanceof Error ? e.message : ''))
        }
      }

      if (notes.length > 0) {
        setDone({ slug: c.url_slug, notes })
        return
      }
      nav(`/c/${c.url_slug}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setBusy(false)
    }
  }

  if (done) {
    return (
      <div className="mx-auto max-w-lg py-12 text-center">
        <p className="text-5xl">🎉</p>
        <p className="mt-4 text-xl font-bold text-gray-900">Ta cagnotte est créée</p>
        {done.notes.map((n) => (
          <p key={n} className="mt-2 text-gray-600">
            {n}
          </p>
        ))}
        <Link to={`/c/${done.slug}`} className="btn btn-primary mt-6">
          Voir ma cagnotte
        </Link>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-lg py-6">
      {/* Progress */}
      <div className="mb-6">
        <div className="mb-2 flex items-center justify-between text-sm text-gray-600">
          <span className="font-semibold text-gray-900">
            Étape {step + 1} sur {total}
          </span>
          {step > 0 && (
            <button
              type="button"
              onClick={() => {
                setError('')
                setStep(step - 1)
              }}
              className="font-medium text-gray-700 hover:text-gray-900"
            >
              ← Retour
            </button>
          )}
        </div>
        <div className="h-2.5 overflow-hidden rounded-full bg-gray-200">
          <div
            className="h-full rounded-full transition-all duration-300"
            style={{ width: `${((step + 1) / total) * 100}%`, backgroundColor: '#0e7c4a' }}
          />
        </div>
      </div>

      <h1 className="text-3xl font-extrabold leading-tight text-gray-900">{HEADINGS[step]}</h1>
      <p className="mt-2 text-gray-600">{HINTS[step]}</p>

      <div className="mt-6 space-y-4">
        {/* Step 1: beneficiary */}
        {step === 0 && (
          <>
            <div className="space-y-2">
              {BENEFICIARY_TYPES.map((b) => (
                <label
                  key={b.value}
                  className={`flex items-center gap-3 rounded-2xl border-2 p-4 text-lg text-gray-900 ${
                    benType === b.value ? 'border-brand-600 bg-brand-50' : 'border-gray-200 bg-white'
                  }`}
                >
                  <input
                    type="radio"
                    name="ben"
                    checked={benType === b.value}
                    onChange={() => setBenType(b.value)}
                    className="h-5 w-5 accent-[#0e7c4a]"
                  />
                  {b.label}
                </label>
              ))}
            </div>
            {benType !== 'self' && (
              <input
                value={benName}
                onChange={(e) => setBenName(e.target.value)}
                maxLength={120}
                placeholder="Nom de la personne ou de la cause"
                className={input}
              />
            )}
            <div>
              <p className="mb-2 font-semibold text-gray-900">Dans quel pays ?</p>
              <select value={country} onChange={(e) => setCountry(e.target.value)} className={input}>
                {UEMOA_COUNTRIES.map((c) => (
                  <option key={c.code} value={c.code}>
                    {c.name}
                  </option>
                ))}
              </select>
            </div>
          </>
        )}

        {/* Step 2: title + goal */}
        {step === 1 && (
          <>
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={120}
              placeholder="Ex. : Opération de ma mère"
              className={input}
            />
            <div>
              <p className="mb-2 font-semibold text-gray-900">
                Objectif <span className="text-sm font-normal text-gray-500">(facultatif)</span>
              </p>
              <div className="relative">
                <input
                  inputMode="numeric"
                  value={goal}
                  onChange={(e) => setGoal(e.target.value.replace(/\D/g, ''))}
                  placeholder="Montant à collecter"
                  className={input + ' pr-20'}
                />
                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500">FCFA</span>
              </div>
              {goalNum >= 1_000_000 && (
                <p className="mt-2 rounded-2xl bg-amber-50 p-3 text-sm text-amber-800">
                  Pour un grand objectif, pense à ajouter aussi un compte bancaire : les comptes mobile money ont des
                  plafonds.
                </p>
              )}
            </div>
          </>
        )}

        {/* Step 3: story + photo */}
        {step === 2 && (
          <>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              maxLength={5000}
              rows={7}
              placeholder="Qui es-tu ? Pour qui collectes-tu ? À quoi servira l'argent ?"
              className={input}
            />
            <div>
              <p className="mb-2 font-semibold text-gray-900">
                Une photo{' '}
                <span className="text-sm font-normal text-gray-500">(facultatif, mais elle aide beaucoup)</span>
              </p>
              <PhotoPicker onChange={setPhoto} />
            </div>
          </>
        )}

        {/* Step 4: payout */}
        {step === 3 && (
          <>
            <div className="grid grid-cols-2 gap-3">
              {PAYOUT_TYPES.map((p) => {
                if (p.value === 'bank') {
                  return (
                    <button
                      type="button"
                      key={p.value}
                      onClick={() => pickPayType('bank')}
                      className={`relative flex h-24 flex-col items-center justify-center gap-2 rounded-2xl border-2 border-dashed bg-gray-50 text-gray-600 ${
                        addBank ? 'border-brand-600' : 'border-gray-300'
                      }`}
                    >
                      <span className="absolute -top-2.5 right-3 rounded-full bg-[#F5B82E] px-2.5 py-0.5 text-[11px] font-extrabold text-gray-900">
                        En plus du mobile money
                      </span>
                      <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#F5B82E] text-gray-900">
                        <Landmark size={22} />
                      </span>
                      <span className="text-sm font-medium">{p.label}</span>
                    </button>
                  )
                }
                return (
                  <button
                    type="button"
                    key={p.value}
                    onClick={() => pickPayType(p.value)}
                    className={`flex h-24 flex-col items-center justify-center gap-2 rounded-2xl border-2 bg-white ${
                      payType === p.value ? 'border-brand-600 ring-2 ring-brand-200' : 'border-gray-200'
                    }`}
                  >
                    {p.logo ? (
                      <img
                        src={p.logo}
                        alt=""
                        className="h-10 w-auto object-contain"
                        onError={(e) => (e.currentTarget.style.display = 'none')}
                      />
                    ) : (
                      <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#F5B82E] text-gray-900">
                        <Landmark size={22} />
                      </span>
                    )}
                    <span className="text-sm font-medium text-gray-900">{p.label}</span>
                  </button>
                )
              })}
            </div>

            {bankTapped && !payType && (
              <p className="rounded-2xl border border-[#F5B82E] bg-[#fdf1cc] p-4 text-sm text-gray-900">
                <b>Choisis d'abord Wave ou Orange Money</b> pour recevoir tes dons. Ensuite, tu pourras ajouter ta
                banque : c'est utile pour les grosses sommes.
              </p>
            )}

            {payType && (
              <>
                <input
                  value={payName}
                  onChange={(e) => setPayName(e.target.value)}
                  maxLength={100}
                  placeholder="Nom du titulaire (comme sur le compte)"
                  className={input}
                />
                <input
                  value={payNumber}
                  onChange={(e) => setPayNumber(e.target.value)}
                  inputMode="tel"
                  placeholder="Numéro, ex. +223 70 00 00 00"
                  className={input}
                />
                <p className="rounded-2xl bg-gray-50 p-3 text-sm text-gray-700">
                  Les dons arrivent directement sur ce compte. Nous ne gardons que nos frais. Tu pourras le changer
                  plus tard, avec une vérification de sécurité.
                </p>

                {/* Bank account, after the mobile money account */}
                <div className="rounded-3xl border-2 border-[#F5B82E] bg-[#fdf1cc] p-5">
                  <div className="flex items-start gap-4">
                    <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-[#F5B82E] text-gray-900">
                      <Landmark size={24} />
                    </span>
                    <div className="min-w-0 flex-1">
                      <p className="text-lg font-extrabold text-gray-900">Tu as aussi un compte bancaire ?</p>
                      <p className="mt-1 text-sm text-gray-700">
                        Ajoute-le : tu pourras recevoir de plus grosses sommes, au-delà des limites du mobile money.
                      </p>
                      {bankRecommended && (
                        <span className="mt-2 inline-block rounded-full bg-brand-600 px-3 py-1 text-xs font-semibold text-white">
                          Recommandé pour ta cagnotte
                        </span>
                      )}
                    </div>
                  </div>

                                   <button
                    type="button"
                    onClick={() => setAddBank(!addBank)}
                    aria-pressed={addBank}
                    className={`mt-4 w-full rounded-full border-2 py-3 font-bold transition active:scale-[0.98] ${
                      addBank
                        ? 'border-brand-600 bg-brand-600 text-white hover:bg-brand-700'
                        : 'border-gray-900 bg-white text-gray-900 hover:bg-gray-100'
                    }`}
                  >
                    {addBank ? '✓ Je veux ajouter ma banque' : 'Ajouter ma banque'}
                  </button>

                  {addBank && (
                    <div className="mt-5 space-y-4">
                      <p className="text-sm text-gray-700">
                        Écris le nom exactement comme sur le document de ta banque. Il doit correspondre à celui de ton
                        compte mobile money.
                      </p>
                      <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-900">Nom de la banque</label>
                        <input
                          value={extraBankName}
                          onChange={(e) => setExtraBankName(e.target.value)}
                          maxLength={100}
                          placeholder="Ex : Coris Bank, BDM, UBA"
                          className={smallInput}
                        />
                      </div>
                      <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-900">Nom du titulaire</label>
                        <input
                          value={extraHolder}
                          onChange={(e) => setExtraHolder(e.target.value)}
                          maxLength={100}
                          placeholder="Nom complet, comme sur le compte"
                          className={smallInput}
                        />
                      </div>
                      <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-900">
                          Numéro de compte (IBAN ou RIB)
                        </label>
                        <input
                          value={extraNumber}
                          onChange={(e) => setExtraNumber(e.target.value)}
                          maxLength={40}
                          placeholder="Ex : ML00 0000 0000 0000 0000"
                          className={smallInput}
                        />
                      </div>
                    </div>
                  )}
                </div>
              </>
            )}
          </>
        )}
      </div>

      {error && <p className="mt-4 text-red-600">{error}</p>}

      <button type="button" onClick={next} disabled={busy} className="btn btn-primary group mt-8 w-full">
        {busy ? (
          'Création…'
        ) : step === total - 1 ? (
          <>
            <Rocket size={22} />
            Lancer ma cagnotte
            <ArrowRight size={22} className="transition-transform group-hover:translate-x-1" />
          </>
        ) : (
          'Continuer'
        )}
      </button>
    </div>
  )
}