export interface SpyCat {
  id?: number;
  name: string;
  years_of_experience: number;
  breed: string;
  salary: number;
  mission_id?: number;
  created_at?: string;
  updated_at?: string;
}

export interface CreateCatRequest {
  name: string;
  years_of_experience: number;
  breed: string;
  salary: number;
}

export interface UpdateSalaryRequest {
  salary: number;
}

export interface CatsResponse {
  cats: SpyCat[];
  breeds: string[];
  total: number;
  limit: number;
  offset: number;
}

export interface BreedsResponse {
  breeds: string[];
}
