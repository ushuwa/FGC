export interface Loan {
  id: number;
  client_id: number;
  client_name: string;
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
  frequency?: string | null;
}

export interface LoansResponse {
  success: boolean;
  message: string;
  data: Loan[];
}

export interface LoanProfileClient {
  id: number;
  first_name: string;
  last_name: string;
  contact_number?: string | null;
  email?: string | null;
  current_address?: string | null;
}

export interface LoanProfileInfo {
  id: number;
  client_id: number;
  pn_number: string;
  loan_type?: string | null;
  principal_amount: number;
  interest_rate: number;
  loan_interest: number;
  pn_value: number;
  loan_term: number;
  frequency?: string | null;
  amortization_amount: number;
  disbursement_date?: string | null;
  maturity_date?: string | null;
  status: string;
 
}

export interface LoanProfileSummary {
  principal_amount: number;
  pn_value: number;
  total_paid: number;
  outstanding_balance: number;
}

export interface LoanAmortization {
  id: number;
  loan_id: number;
  due_date: string;
  principal_amount: number;
  interest_amount: number;
  total_amount: number;
  paid_principal_amount: number;
  paid_interest_amount: number;
  status: string;
}

export interface LoanPayment {
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

export interface LoanProfileResponse {
  success: boolean;
  message: string;
  data: {
    loan: LoanProfileInfo;
    client: LoanProfileClient;
    summary: LoanProfileSummary;
    amortizations: LoanAmortization[];
    payments: LoanPayment[];
  };
}


export interface CreatePaymentRequest {
  payment_date: string;
  amount_paid: number;
  payment_channel?: string | null;
  reference_number?: string | null;
}

export interface CreateLoanRequest {
  client_id: number;
  pn_number: string;
  loan_type?: string | null;
  principal_amount: number;
  interest_rate: number;
  loan_interest: number;
  pn_value: number;
  loan_term: number;
  amortization_amount: number;
  disbursement_date: string;
  maturity_date?: string | null;
  frequency?: string | null;
  status: string;
}

export type UpdateLoanRequest = CreateLoanRequest;