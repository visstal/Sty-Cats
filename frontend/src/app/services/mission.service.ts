import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Mission, CreateMissionRequest } from '../interfaces/mission.interface';

@Injectable({
  providedIn: 'root'
})
export class MissionService {
  private apiUrl = 'http://localhost:3001/api/v1/agency/missions';

  constructor(private http: HttpClient) { }

  // Get all missions
  getMissions(): Observable<Mission[]> {
    return this.http.get<Mission[]>(this.apiUrl);
  }

  // Get a specific mission by ID
  getMission(id: number): Observable<Mission> {
    return this.http.get<Mission>(`${this.apiUrl}/${id}`);
  }

  // Create a new mission
  createMission(mission: CreateMissionRequest): Observable<Mission> {
    return this.http.post<Mission>(this.apiUrl, mission);
  }

  // Delete a mission
  deleteMission(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  // Assign a cat to a mission
  assignCatToMission(missionId: number, catId: number): Observable<Mission> {
    return this.http.post<Mission>(`${this.apiUrl}/${missionId}/assign`, { cat_id: catId });
  }

  // Get free cats available for assignment
  getFreeCats(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/free-cats`);
  }

  // Add a target to a mission
  addTargetToMission(missionId: number, target: { name: string; country: string; notes?: string }): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${missionId}/targets`, target);
  }

  // Delete a target from a mission (only if status is 'init')
  deleteTargetFromMission(missionId: number, targetId: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${missionId}/targets/${targetId}`);
  }
}
