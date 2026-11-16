export interface Campaign {
  id: string
  organizer_id: string
  title: string
  description?: string
  goal_amount?: number
  current_amount: number
  currency: string
  beneficiary_type: string
  beneficiary_name?: string
  country: string
  url_slug: string
  status: string
  created_at: string
  cover_url?: string // the API does not send it yet (see note below)
    cover_sensitive?: boolean
}

export interface CampaignList {
  campaigns: Campaign[]
}

export interface PublicDonation {
  id: string
  donor_name?: string
  amount: number
  comment?: string
  reply?: string
  reply_at?: string
  created_at: string
}

export interface CampaignDetail extends Campaign {
  description?: string
  organizer_name?: string
}

export interface DonationCreated {
  id: string
  campaign_id: string
  amount: number
  fee_amount: number
  tip_amount: number
  total_charged: number
  currency: string
  status: string
}

export interface PayoutAccount {
  id: string
  type: string
  account_name: string
  account_number: string
  bank_name?: string
  bank_code?: string
  country: string
  is_active: boolean
  verified_at?: string
  created_at: string
}

export interface OrganizerDonation {
  id: string
  donor_name: string
  is_anonymous: boolean
  amount: number
  currency: string
  comment?: string
  comment_hidden: boolean
  created_at: string
  reply?: string
  reply_at?: string
}