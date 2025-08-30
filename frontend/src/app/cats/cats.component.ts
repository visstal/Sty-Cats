import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CatService } from './cat.service';
import { SpyCat, CreateCatRequest, CatsResponse, BreedsResponse } from './cat.interface';

@Component({
  selector: 'app-cats',
  template: `
    <div class="cats-container">
      <!-- Add Cat Form -->
      <div class="spy-card">
        <h2>🕵️ Add New Spy Cat Agent</h2>
        
        <form [formGroup]="catForm" (ngSubmit)="onSubmit()" class="spy-form">
          <div class="form-group">
            <label class="spy-label">Agent Name</label>
            <input 
              class="spy-input" 
              type="text" 
              formControlName="name"
              placeholder="Enter spy cat name"
              maxlength="100">
            <div *ngIf="catForm.get('name')?.invalid && catForm.get('name')?.touched" class="validation-error">
              <span *ngIf="catForm.get('name')?.errors?.['required']">🚨 Agent name is required</span>
              <span *ngIf="catForm.get('name')?.errors?.['minlength']">🚨 Agent name must be at least 2 characters</span>
              <span *ngIf="catForm.get('name')?.errors?.['maxlength']">🚨 Agent name must not exceed 100 characters</span>
            </div>
          </div>
          
          <div class="form-group">
            <label class="spy-label">Breed</label>
            <select 
              class="spy-input" 
              formControlName="breed">
              <option value="" disabled>Select cat breed</option>
              <option *ngFor="let breed of breeds" [value]="breed">{{breed}}</option>
            </select>
            <div *ngIf="catForm.get('breed')?.invalid && catForm.get('breed')?.touched" class="validation-error">
              <span *ngIf="catForm.get('breed')?.errors?.['required']">🚨 Breed selection is required</span>
            </div>
          </div>
          
          <div class="form-group">
            <label class="spy-label">Years of Experience</label>
            <input 
              class="spy-input" 
              type="number" 
              formControlName="years_of_experience"
              placeholder="Years in the field"
              min="0"
              max="50">
            <div *ngIf="catForm.get('years_of_experience')?.invalid && catForm.get('years_of_experience')?.touched" class="validation-error">
              <span *ngIf="catForm.get('years_of_experience')?.errors?.['required']">🚨 Experience is required</span>
              <span *ngIf="catForm.get('years_of_experience')?.errors?.['min']">🚨 Experience must be at least 0 years</span>
              <span *ngIf="catForm.get('years_of_experience')?.errors?.['max']">🚨 Experience cannot exceed 50 years</span>
            </div>
          </div>
          
          <div class="form-group">
            <label class="spy-label">Salary ($)</label>
            <input 
              class="spy-input" 
              type="number" 
              formControlName="salary"
              placeholder="Annual salary"
              min="0">
            <div *ngIf="catForm.get('salary')?.invalid && catForm.get('salary')?.touched" class="validation-error">
              <span *ngIf="catForm.get('salary')?.errors?.['required']">🚨 Salary is required</span>
              <span *ngIf="catForm.get('salary')?.errors?.['min']">🚨 Salary must be at least $0</span>
            </div>
          </div>
          
          <div class="form-group full-width">
            <button 
              type="submit" 
              class="spy-btn success"
              [disabled]="isLoading">
              {{isLoading ? '🔄 Adding Agent...' : '➕ Add Spy Cat'}}
            </button>
          </div>
        </form>
      </div>

      <!-- Success/Error Messages -->
      <div *ngIf="message" class="alert" [class]="messageType">
        {{message}}
      </div>

      <!-- Cats List -->
      <div class="spy-card">
        <h2>🐱 Active Spy Cat Agents</h2>
        
        <div *ngIf="isLoading && cats.length === 0" class="loading">
          <div class="loading-spinner"></div>
        </div>
        
        <div *ngIf="!isLoading && cats.length === 0" class="no-data">
          <p>No spy cats recruited yet. Add your first agent! 🕵️‍♀️</p>
        </div>
        
        <table *ngIf="cats.length > 0" class="spy-table">
          <thead>
            <tr>
              <th>🆔 ID</th>
              <th>🐱 Name</th>
              <th>🏷️ Breed</th>
              <th>⏱️ Experience</th>
              <th>💰 Salary</th>
              <th>🎯 Status</th>
              <th>📅 Recruited</th>
              <th>🔧 Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr *ngFor="let cat of cats; trackBy: trackByCatId">
              <td>{{cat.id}}</td>
              <td>{{cat.name}}</td>
              <td>{{cat.breed}}</td>
              <td>{{cat.years_of_experience}} years</td>
              <td *ngIf="editingSalary !== cat.id">
                \${{cat.salary}}
                <button 
                  class="spy-btn" 
                  style="margin-left: 0.5rem; padding: 0.3rem 0.8rem; font-size: 0.8rem;"
                  (click)="startEditSalary(cat)">
                  ✏️
                </button>
              </td>
              <td *ngIf="editingSalary === cat.id">
                <div style="display: flex; gap: 0.5rem; align-items: center;">
                  <input 
                    class="spy-input" 
                    type="number" 
                    [(ngModel)]="newSalary"
                    style="width: 100px; padding: 0.3rem;">
                  <button 
                    class="spy-btn success" 
                    style="padding: 0.3rem 0.8rem; font-size: 0.8rem;"
                    (click)="updateSalary(cat.id!)">
                    ✅
                  </button>
                  <button 
                    class="spy-btn" 
                    style="padding: 0.3rem 0.8rem; font-size: 0.8rem;"
                    (click)="cancelEditSalary()">
                    ❌
                  </button>
                </div>
              </td>
              <td>
                <span 
                  class="status-badge"
                  [class.assigned]="isAssignedToMission(cat)"
                  [class.standby]="!isAssignedToMission(cat)">
                  {{isAssignedToMission(cat) ? '🎯 ASSIGNED' : '🟡 STANDBY'}}
                </span>
              </td>
              <td>{{cat.created_at | date:'short'}}</td>
              <td>
                <button 
                  class="spy-btn danger" 
                  style="padding: 0.5rem 1rem; font-size: 0.9rem;"
                  (click)="deleteCat(cat.id!)"
                  [disabled]="isLoading || isAssignedToMission(cat)"
                  [title]="isAssignedToMission(cat) ? 'Cannot delete cat assigned to mission' : 'Terminate spy cat'">
                  🗑️ Terminate
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  `,
  styleUrls: ['./cats.component.css']
})
export class CatsComponent implements OnInit {
  catForm: FormGroup;
  cats: SpyCat[] = [];
  breeds: string[] = [];
  isLoading = false;
  message = '';
  messageType = '';
  editingSalary: number | null = null;
  newSalary = 0;

  constructor(
    private fb: FormBuilder,
    private catService: CatService
  ) {
    this.catForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(1), Validators.maxLength(100)]],
      breed: ['', [Validators.required, Validators.minLength(1), Validators.maxLength(100)]],
      years_of_experience: [0, [Validators.required, Validators.min(0), Validators.max(50)]],
      salary: [0, [Validators.required, Validators.min(0)]]
    });
  }

  ngOnInit() {
    this.loadCats();
    this.loadBreeds();
  }

  loadCats() {
    this.isLoading = true;
    this.catService.getAllCats().subscribe({
      next: (response: CatsResponse) => {
        this.cats = response.cats || [];
        this.breeds = response.breeds || [];
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading cats:', error);
        this.showMessage('Failed to load spy cats. Check if the backend is running.', 'error');
        this.isLoading = false;
      }
    });
  }

  loadBreeds() {
    this.catService.getBreeds().subscribe({
      next: (response: BreedsResponse) => {
        this.breeds = response.breeds || [];
      },
      error: (error) => {
        console.error('Error loading breeds:', error);
      }
    });
  }

  onSubmit() {
    // Always allow submission to get detailed backend validation
    this.isLoading = true;
    const catData: CreateCatRequest = this.catForm.value;
    
    this.catService.createCat(catData).subscribe({
      next: (newCat) => {
        this.cats.unshift(newCat);
        this.catForm.reset();
        this.showMessage('✅ Spy cat "' + newCat.name + '" has been recruited successfully!', 'success');
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error creating cat:', error);
        let errorMessage = '❌ Failed to recruit spy cat. Please try again.';
        
        // Extract detailed validation messages from backend response
        if (error.status === 400 && error.error) {
          if (error.error.details) {
            errorMessage = `❌ ${error.error.details}`;
          } else if (error.error.error) {
            errorMessage = `❌ ${error.error.error}`;
          }
        } else if (error.error && error.error.details) {
          errorMessage = `❌ ${error.error.details}`;
        }
        
        this.showMessage(errorMessage, 'error');
        this.isLoading = false;
      }
    });
  }

  startEditSalary(cat: SpyCat) {
    this.editingSalary = cat.id!;
    this.newSalary = cat.salary;
  }

  cancelEditSalary() {
    this.editingSalary = null;
    this.newSalary = 0;
  }

  updateSalary(catId: number) {
    if (this.newSalary < 0) return;
    
    this.isLoading = true;
    this.catService.updateCatSalary(catId, { salary: this.newSalary }).subscribe({
      next: (updatedCat) => {
        const index = this.cats.findIndex(c => c.id === catId);
        if (index !== -1) {
          this.cats[index] = updatedCat;
        }
        this.editingSalary = null;
        this.showMessage('💰 Salary updated successfully!', 'success');
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error updating salary:', error);
        let errorMessage = '❌ Failed to update salary. Please try again.';
        
        // Extract detailed validation messages from backend response
        if (error.status === 400 && error.error) {
          if (error.error.details) {
            errorMessage = `❌ ${error.error.details}`;
          } else if (error.error.error) {
            errorMessage = `❌ ${error.error.error}`;
          }
        } else if (error.error && error.error.details) {
          errorMessage = `❌ ${error.error.details}`;
        }
        
        this.showMessage(errorMessage, 'error');
        this.isLoading = false;
      }
    });
  }

  deleteCat(catId: number) {
    const cat = this.cats.find(c => c.id === catId);
    if (cat && confirm('Are you sure you want to terminate spy cat "' + cat.name + '"?')) {
      this.isLoading = true;
      
      this.catService.deleteCat(catId).subscribe({
        next: () => {
          this.cats = this.cats.filter(c => c.id !== catId);
          this.showMessage('🗑️ Spy cat "' + cat.name + '" has been terminated.', 'success');
          this.isLoading = false;
        },
        error: (error) => {
          console.error('Error deleting cat:', error);
          let errorMessage = '❌ Failed to terminate spy cat. Please try again.';
          
          // Check if it's a conflict error (cat assigned to mission)
          if (error.status === 409 && error.error && error.error.details) {
            errorMessage = `❌ ${error.error.details}`;
          } else if (error.error && error.error.details) {
            errorMessage = `❌ ${error.error.details}`;
          }
          
          this.showMessage(errorMessage, 'error');
          this.isLoading = false;
        }
      });
    }
  }

  trackByCatId(index: number, cat: SpyCat): number {
    return cat.id || index;
  }

  private showMessage(message: string, type: string) {
    this.message = message;
    this.messageType = type;
    setTimeout(() => {
      this.message = '';
      this.messageType = '';
    }, 5000);
  }

  isAssignedToMission(cat: SpyCat): boolean {
    return !!(cat.mission_id);
  }

  getMissionStatus(cat: SpyCat): string {
    return cat.mission_id ? 'ASSIGNED' : 'STANDBY';
  }
}
