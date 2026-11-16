// Same formula as the backend: (amount * 6 + 50) / 100, integer. The server is the source of truth.
export const PLATFORM_FEE_PERCENT = 6
export const calcFee = (amount: number): number =>
  Math.floor((amount * PLATFORM_FEE_PERCENT + 50) / 100)

export interface PayMethod {
  value: string // must match models.PaymentMethod in the backend
  label: string
  logo: string
  max: number // per-transaction limit in FCFA, 0 = no limit
}

export const PAY_METHODS: PayMethod[] = [
  { value: 'wave', label: 'Wave', logo: '/logos/wave.png', max: 0 },
  { value: 'orange_money', label: 'Orange Money', logo: '/logos/orange-money.png', max: 1_000_000 },
  { value: 'moov', label: 'Moov Money', logo: '/logos/moov-money.png', max: 500_000 },
  { value: 'card', label: 'Carte Visa / Mastercard', logo: '/logos/visa.png', max: 0 },
]

export const UEMOA_COUNTRIES = [
  { code: 'ML', name: 'Mali' },
  { code: 'SN', name: 'Sénégal' },
  { code: 'CI', name: "Côte d'Ivoire" },
  { code: 'BF', name: 'Burkina Faso' },
  { code: 'BJ', name: 'Bénin' },
  { code: 'TG', name: 'Togo' },
  { code: 'NE', name: 'Niger' },
  
]

export const BENEFICIARY_TYPES = [
  { value: 'self', label: 'Pour moi' },
  { value: 'other_person', label: 'Pour une autre personne' },
  { value: 'cause_collective', label: 'Pour une cause collective' },
]

export const PAYOUT_TYPES = [
  { value: 'wave', label: 'Wave', logo: '/logos/wave.png' },
  { value: 'orange_money', label: 'Orange Money', logo: '/logos/orange-money.png' },
  { value: 'moov_money', label: 'Moov Money', logo: '/logos/moov-money.png' },
  { value: 'bank', label: 'Compte bancaire', logo: '' },
]