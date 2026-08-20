import api from "./api";

export interface PARSummary {
  par_loans: number;
  default_amount: number;
  par_ratio: number;
}

export interface PARLoan {
  id: number;
  pn_number: string;
  client_name: string;
  due_date: string;
  days_past_due: number;
  default_amount: number;
  status: string;
}

export interface PARPagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PARResponse {
  summary: PARSummary;
  loans: PARLoan[];
  pagination: PARPagination;
  aging: PARAging[];
}

interface APIResponse {
  success: boolean;
  data: PARResponse;
}

export interface PARFilters {
  search?: string;
  status?: string;
  aging?: string;
  page?: number;
  limit?: number;
}

export interface PARAging{
    aging: "1-30" | "31-60" | "61-90" | "90+";
    loans: number;
    default_amount: number;
}

export const getPortfolioAtRisk =
  async (
    filters: PARFilters = {},
  ): Promise<PARResponse> => {

    const params = new URLSearchParams();

    if (filters.search) {
      params.set(
        "search",
        filters.search,
      );
    }

    if (
      filters.status &&
      filters.status !== "ALL"
    ) {
      params.set(
        "status",
        filters.status,
      );
    }

    if (
      filters.aging &&
      filters.aging !== "ALL"
    ) {
      params.set(
        "aging",
        filters.aging,
      );
    }

    params.set(
      "page",
      String(filters.page ?? 1),
    );

    params.set(
      "limit",
      String(filters.limit ?? 10),
    );

    const response =
      await api.get<APIResponse>(
        `/portfolio-at-risk/?${params.toString()}`,
      );

    return response.data.data;
  };