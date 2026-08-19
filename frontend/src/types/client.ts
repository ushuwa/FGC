export interface Client {
  id: number;
  first_name: string;
  last_name: string;
  contact_number?: string | null;
  email?: string | null;
  current_address?: string | null;
  created_at: string;
}

export interface CreateClientRequest {
  first_name: string;
  last_name: string;
  contact_number?: string;
  email?: string;
  current_address?: string;
}

export interface UpdateClientRequest {
  first_name: string;
  last_name: string;
  contact_number?: string;
  email?: string;
  current_address?: string;
}

export interface ClientProfile {
  id: number;
  first_name: string;
  last_name: string;
  contact_number?: string | null;
  email?: string | null;
  current_address?: string | null;
}

export interface ClientSummary {
  total_loans: number;
  active_loans: number;
  total_principal: number;
  total_paid: number;
  total_outstanding: number;
}

export interface ClientLoanSummary {
  id: number;
  pn_number: string;
  loan_type?: string | null;
  principal_amount: number;
  interest_rate: number;
  loan_interest: number;
  pn_value: number;
  loan_term: number;
  amortization_amount: number;
  disbursement_date?: string | null;
  maturity_date?: string | null;
  status: string;
  total_paid: number;
  outstanding_balance: number;
}

export interface ClientPayment {
  id: number;
  loan_id: number;
  payment_date: string;
  amount_paid: number;
  payment_channel?: string | null;
  reference_number?: string | null;
  principal_applied: number;
  interest_applied: number;
  outstanding_balance: number;
}

export interface ClientProfileResponse {
  client: ClientProfile;
  summary: ClientSummary;
  loans: ClientLoanSummary[];
  payments: ClientPayment[];
}