export interface Cat {
  id: number;
  name: string;
  breed: string;
  salary: number;
  mission_years: number;
  created_at: string;
  updated_at: string;
}

export interface CreateCatRequest {
  name: string;
  breed: string;
  salary: number;
  mission_years: number;
}

export interface UpdateCatRequest {
  name?: string;
  breed?: string;
  salary?: number;
  mission_years?: number;
}
