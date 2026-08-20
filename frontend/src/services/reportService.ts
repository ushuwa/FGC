import api from "./api";

export async function downloadCollectionReport(
  dateFrom: string,
  dateTo: string,
): Promise<Blob> {

  const response =
    await api.get(
      "/reports/summary-collection/pdf",
      {
        params: {
          date_from: dateFrom,
          date_to: dateTo,
        },
        responseType: "blob",
      },
    );

  return response.data;
}

export async function downloadAmortizationReport(
  loanId: number,
): Promise<Blob> {

  const response = await api.get(
    `/reports/amortization/${loanId}/pdf`,
    {
      responseType: "blob",
    },
  );

  return response.data;
}

export async function downloadPARReport(): Promise<Blob> {

  const response = await api.get(
    "/reports/portfolio-at-risk/pdf",
    {
      responseType: "blob",
    },
  );

  return response.data;
}

export async function downloadLoanPortfolioReport(): Promise<Blob> {

  const response = await api.get(
    "/reports/loan-portfolio/pdf",
    {
      responseType: "blob",
    },
  );

  return response.data;
}

export async function downloadLoanMaturityReport(): Promise<Blob> {

  const response = await api.get(
    "/reports/loan-maturity/pdf",
    {
      responseType: "blob",
    },
  );

  return response.data;
}