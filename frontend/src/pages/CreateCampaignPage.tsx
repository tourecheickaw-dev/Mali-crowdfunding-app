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

  // main payout account
  const [payType, setPayType] = useState('')
  const [payName, setPayName] = useState('')
  const [payBankName, setPayBankName] = useState('')
  const [payNumber, setPayNumber] = useState('')

  // optional backup bank account (only when the main one is mobile money)
  const [addBank, setAddBank] = useState(false)
  const [extraBankName, setExtraBankName] = useState('')
  const [extraHolder, setExtraHolder] = useState('')
  const [extraNumber, setExtraNumber] = useState('')

  const [payoutCreated, setPayoutCreated] = useState(false)
  const [bankCreated, setBankCreated] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [photoFailedSlug, setPhotoFailedSlug] = useState('')

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
  const isBank = payType === 'bank'
  const cleanNumber = payNumber.replace(/[\s.-]/g, '')
  const cleanExtra = extraNumber.replace(/\s/g, '')

  function validate(s: number): string {
    if (s === 0 && benType !== 'self' && !benName.trim()) return 'Écris le nom de la personne ou de la cause.'
    if (s === 1 && !title.trim()) return 'Donne un titre à ta cagnotte.'
    if (s === 2 && description.trim().length < 20)
      return 'Raconte ton histoire en quelques phrases (20 caractères minimum).'
    if (s === 3) {
      if (!payType) return "Choisis où recevoir l'argent."
      if (!payName.trim()) return 'Écris le nom du titulaire du compte.'
      if (isBank) {
        if (!payBankName.trim()) return 'Écris le nom de la banque.'
        if (cleanNumber.length < 8) return 'Numéro de compte trop court.'
      } else {
        if (!PHONE_RE.test(cleanNumber)) return 'Numéro au format international, par exemple +223 70 00 00 00.'
        if (addBank) {
          if (!extraBankName.trim()) return 'Écris le nom de la banque.'
          if (!extraHolder.trim()) return 'Écris le nom du titulaire du compte bancaire.'
          if (cleanExtra.length < 8) return 'Numéro du compte bancaire trop court.'
        }
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
            bank_name: isBank ? payBankName.trim() : '',
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

      if (!hasPayout && addBank && !isBank && !bankCreated) {
        // backup account: created but NOT activated, donations keep going to the main account
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

      if (photo) {
        try {
          await attachCover(c.id, photo.url, photo.sensitive)
        } catch {
          setPhotoFailedSlug(c.url_slug)
          return
        }
      }
      nav(`/c/${c.url_slug}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setBusy(false)
    }
  }

  if (photoFailedSlug) {
    return (
      <div className="mx-auto max-w-lg py-12 text-center">
        <p className="text-5xl">🎉</p>
        <p className="mt-4 text-xl font-bold text-gray-900">Ta cagnotte est créée</p>
        <p className="mt-2 text-gray-600">Mais la photo n'a pas pu être ajoutée. Tu pourras réessayer plus tard.</p>
        <Link to={`/c/${photoFailedSlug}`} className="btn btn-primary mt-6">
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
              {PAYOUT_TYPES.map((p) => (
                <button
                  type="button"
                  key={p.value}
                  onClick={() => setPayType(p.value)}
                  className={`flex h-24 flex-col items-center justify-center gap-2 rounded-2xl border-2 bg-white ${
                    payType === p.value ? 'border-brand-600 ring-2 ring-brand-200' : 'border-gray-200'
                  }`}
                >
                  {p.value === 'bank' || !p.logo ? (
                    <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#F5B82E] text-gray-900">
                      <Landmark size={22} />
                    </span>
                  ) : (
                    <img
                      src={p.logo}
                      alt=""
                      className="h-10 w-auto object-contain"
                      onError={(e) => (e.currentTarget.style.display = 'none')}
                    />
                  )}
                  <span className="text-sm font-medium text-gray-900">{p.label}</span>
                </button>
              ))}
            </div>

            {payType && (
              <>
                <input
                  value={payName}
                  onChange={(e) => setPayName(e.target.value)}
                  maxLength={100}
                  placeholder="Nom du titulaire (comme sur le compte)"
                  className={input}
                />
                {isBank && (
                  <input
                    value={payBankName}
                    onChange={(e) => setPayBankName(e.target.value)}
                    maxLength={100}
                    placeholder="Nom de la banque"
                    className={input}
                  />
                )}
                <input
                  value={payNumber}
                  onChange={(e) => setPayNumber(e.target.value)}
                  inputMode={isBank ? 'text' : 'tel'}
                  placeholder={isBank ? 'Numéro de compte (IBAN ou RIB)' : 'Numéro, ex. +223 70 00 00 00'}
                  className={input}
                />
                              <p className="rounded-2xl bg-gray-50 p-3 text-sm text-gray-700">
                  {isBank
                    ? 'Ce compte sera vérifié avant d’être activé. Tu pourras le changer plus tard, avec une vérification de sécurité.'
                    : 'Les dons arrivent directement sur ce compte. Nous ne gardons que nos frais. Tu pourras le changer plus tard, avec une vérification de sécurité.'}
                </p>
              </>
            )}

            {/* Optional backup bank account */}
            {payType && !isBank && (
              <div
                className={`mt-4 rounded-3xl border-2 p-5 transition ${
                  addBank ? 'border-brand-600 bg-brand-50' : 'border-gray-200 bg-white'
                }`}
              >
                <div className="flex items-start gap-4">
                  <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-[#F5B82E] text-gray-900">
                    <Landmark size={24} />
                  </span>
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <p className="text-lg font-extrabold text-gray-900">Ajouter un compte bancaire</p>
                      {bankRecommended && (
                        <span className="rounded-full bg-brand-600 px-3 py-1 text-xs font-semibold text-white">
                          Recommandé
                        </span>
                      )}
                    </div>
                    <p className="mt-1 text-sm text-gray-600">
                      Facultatif. Les comptes mobile money ont des plafonds : si le tien est plein, tu pourras passer
                      sur ta banque.
                    </p>
                  </div>
                  <button
                    type="button"
                    role="switch"
                    aria-checked={addBank}
                    aria-label="Ajouter un compte bancaire"
                    onClick={() => setAddBank(!addBank)}
                    className="relative h-7 w-12 shrink-0 rounded-full transition"
                    style={{ backgroundColor: addBank ? '#0e7c4a' : '#d1d5db' }}
                  >
                    <span
                      className={`absolute top-1 h-5 w-5 rounded-full bg-white shadow transition-all ${
                        addBank ? 'left-6' : 'left-1'
                      }`}
                    />
                  </button>
                </div>

                {addBank && (
                  <div className="mt-5 space-y-4">
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
