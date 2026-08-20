import api from "./api";

export interface DashboardSummary {
  total_clients: number;
  total_loans: number;
  active_loans: number;
  paid_loans: number;
  total_principal: number;
  total_pn_value: number;
  total_collected: number;
  total_outstanding: number;
}

interface DashboardSummaryResponse {
  success: boolean;
  data: DashboardSummary;
}

export const getDashboardSummary =
  async (): Promise<DashboardSummary> => {
    const response =
      await api.get<DashboardSummaryResponse>(
        "/dashboard/summary",
      );

    return response.data.data;
  };