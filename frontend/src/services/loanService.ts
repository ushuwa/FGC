import api from "./api";

import type {
  Loan,
  LoansResponse,
  LoanProfileResponse,
  LoanPayment,
  CreatePaymentRequest,
  CreateLoanRequest,
  UpdateLoanRequest,
} from "../types/loan";

export async function getLoans(
  search = "",
  status = "",
): Promise<Loan[]> {
  const response =
    await api.get<LoansResponse>(
      "/loans",
      {
        params: {
          search,
          status,
        },
      },
    );

  return response.data.data;
}

export async function getLoanProfile(
  id: number,
) {
  const response =
    await api.get<LoanProfileResponse>(
      `/loans/${id}`,
    );

  return response.data.data;
}

export async function getLoanPayments(
  id: number,
) {
  const response =
    await api.get<{
      success: boolean;
      message: string;
      data: LoanPayment[];
    }>(
      `/loans/${id}/payments`,
    );

  return response.data.data;
}

export async function createLoanPayment(
  loanId: number,
  payload: CreatePaymentRequest,
) {
  const response =
    await api.post<{
      success: boolean;
      message: string;
      data: LoanPayment;
    }>(
      `/loans/${loanId}/payments`,
      payload,
    );

  return response.data.data;
}

export async function rebuildAmortization(
  loanId: number,
): Promise<void> {
  await api.post(
    `/loans/${loanId}/rebuild-amortization`,
  );
}

export async function createLoan(
  payload: CreateLoanRequest,
) {
  const response =
    await api.post<{
      success: boolean;
      message: string;
      data: Loan;
    }>(
      "/loans",
      payload,
    );

  return response.data.data;
}

export async function updateLoan(
  id: number,
  payload: UpdateLoanRequest,
) {
  const response =
    await api.put<{
      success: boolean;
      message: string;
      data: Loan;
    }>(
      `/loans/${id}`,
      payload,
    );

  return response.data.data;
}