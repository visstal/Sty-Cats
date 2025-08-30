import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { SpyCat, CreateCatRequest, UpdateSalaryRequest, CatsResponse, BreedsResponse } from './cat.interface';

@Injectable({
  providedIn: 'root'
})
export class CatService {
  private commonApiUrl = 'http://localhost:3001/api/v1/cats';  // Common endpoints for both modes
  private agencyApiUrl = 'http://localhost:3001/api/v1/agency/cats';  // Agency-only endpoints

  constructor(private http: HttpClient) {}

  // Common endpoints (used by both Agency and Spy Cat modes)
  getAllCats(): Observable<CatsResponse> {
    return this.http.get<CatsResponse>(this.commonApiUrl);
  }

  getCat(id: number): Observable<SpyCat> {
    return this.http.get<SpyCat>(`${this.commonApiUrl}/${id}`);
  }

  getBreeds(): Observable<BreedsResponse> {
    return this.http.get<BreedsResponse>(`${this.commonApiUrl}/breeds`);
  }

  // Agency-only endpoints (administrative operations)
  createCat(cat: CreateCatRequest): Observable<SpyCat> {
    return this.http.post<SpyCat>(this.agencyApiUrl, cat);
  }

  updateCatSalary(id: number, request: UpdateSalaryRequest): Observable<SpyCat> {
    return this.http.put<SpyCat>(`${this.agencyApiUrl}/${id}/salary`, request);
  }

  deleteCat(id: number): Observable<void> {
    return this.http.delete<void>(`${this.agencyApiUrl}/${id}`);
  }
}
