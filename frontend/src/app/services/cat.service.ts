import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Cat, CreateCatRequest, UpdateCatRequest } from '../interfaces/cat.interface';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class CatService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getCats(): Observable<Cat[]> {
    return this.http.get<Cat[]>(`${this.apiUrl}/cats`);
  }

  getCat(id: number): Observable<Cat> {
    return this.http.get<Cat>(`${this.apiUrl}/cats/${id}`);
  }

  createCat(catData: CreateCatRequest): Observable<Cat> {
    return this.http.post<Cat>(`${this.apiUrl}/cats`, catData);
  }

  updateCat(id: number, catData: UpdateCatRequest): Observable<Cat> {
    return this.http.put<Cat>(`${this.apiUrl}/cats/${id}`, catData);
  }

  deleteCat(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/cats/${id}`);
  }
}
