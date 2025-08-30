export interface Target {
  id?: number;
  mission_id?: number;
  name: string;
  country: string;
  notes?: string;
  status: string;
  created_at?: Date | string;
  updated_at?: Date | string;
}

export interface CreateTargetRequest {
  name: string;
  country: string;
}

export interface Mission {
  id?: number;
  name: string;
  description: string;
  start_date?: Date | string | null;
  end_date?: Date | string | null;
  cat_id?: number;
  is_completed?: boolean;
  completed_at?: Date | string | null;
  created_at?: Date | string;
  updated_at?: Date | string;
  targets?: Target[];
  cat?: {
    id: number;
    name: string;
    years_of_experience: number;
    breed: string;
    salary: number;
    created_at: Date | string;
    updated_at: Date | string;
  };
}

export interface CreateMissionRequest {
  name: string;
  description: string;
  start_date?: Date | string;
  end_date?: Date | string;
  targets: CreateTargetRequest[]; // Now required (min 1, max 3)
}
