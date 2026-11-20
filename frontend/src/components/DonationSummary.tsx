import { fcfa } from '../api'

export const FEE_PERCENT = 6

export function computeFee(amount: number): number {
  if (!Number.isFinite(amount) || amount <= 0) return 0
  return Math.floor((amount * FEE_PERCENT + 50) / 100)
}

export default function DonationSummary({ amount }: { amount: number }) {
  const fee = computeFee(amount)
  const total = amount + fee

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5">
      <div className="flex items-center justify-between text-gray-900">
        <span className="font-semibold">Ton don</span>
        <span className="font-extrabold">{fcfa(amount)}</span>
      </div>
      <p className="mt-1 text-sm text-gray-600">100 % de ce montant va à l'organisateur.</p>

      <div className="mt-4 flex items-center justify-between text-gray-900">
        <span className="font-semibold">Frais de service Demso ({FEE_PERCENT} %)</span>
        <span className="font-bold">{fcfa(fee)}</span>
      </div>
      <p className="mt-1 text-sm text-gray-600">Paiement sécurisé et fonctionnement de Demso.</p>

      <div className="mt-4 flex items-center justify-between border-t-2 border-gray-900 pt-4">
        <span className="text-lg font-extrabold text-gray-900">Total à payer</span>
        <span className="text-2xl font-extrabold text-brand-700">{fcfa(total)}</span>
      </div>
    </div>
  )
}