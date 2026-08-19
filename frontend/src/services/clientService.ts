import api from "./api";

import type {
  Client,
  CreateClientRequest,
  UpdateClientRequest,
  ClientProfileResponse,
} from "../types/client";

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}



export async function getClientProfile(
  id: number,
): Promise<ClientProfileResponse> {

  const response =
    await api.get<
      ApiResponse<ClientProfileResponse>
    >(
      `/clients/${id}/profile`,
    );

  return response.data.data;
}
export async function getClients(
  search = "",
): Promise<Client[]> {
  const response =
    await api.get<ApiResponse<Client[]>>(
      "/clients",
      {
        params: {
          search,
        },
      },
    );

  return response.data.data;
}

export async function getClient(
  id: number,
): Promise<Client> {
  const response =
    await api.get<ApiResponse<Client>>(
      `/clients/${id}`,
    );

  return response.data.data;
}

export async function createClient(
  data: CreateClientRequest,
): Promise<Client> {
  const response =
    await api.post<ApiResponse<Client>>(
      "/clients",
      data,
    );

  return response.data.data;
}

export async function updateClient(
  id: number,
  data: UpdateClientRequest,
): Promise<Client> {
  const response =
    await api.put<ApiResponse<Client>>(
      `/clients/${id}`,
      data,
    );

  return response.data.data;
}

export async function deleteClient(
  id: number,
): Promise<void> {
  await api.delete(
    `/clients/${id}`,
  );
}


