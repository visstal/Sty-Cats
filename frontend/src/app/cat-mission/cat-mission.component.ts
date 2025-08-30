import { Component, Input, OnInit, OnChanges, SimpleChanges } from '@angular/core';
import { Mission, Target } from '../interfaces/mission.interface';

interface SpyCat {
  id: number;
  name: string;
  breed: string;
  yearsOfExperience: number;
  salary: number;
}

@Component({
  selector: 'app-cat-mission',
  template: `
    <div class="cat-mission-container" *ngIf="selectedCat">
      <div class="cat-header">
        <h2>🕵️ Agent {{ selectedCat.name }} - Mission Dashboard</h2>
        <div class="cat-info">
          <span class="cat-detail">🏷️ {{ selectedCat.breed }}</span>
          <span class="cat-detail">⏱️ {{ selectedCat.yearsOfExperience }} years experience</span>
          <span class="cat-detail">💰 {{ selectedCat.salary | currency }}</span>
        </div>
      </div>

      <div *ngIf="loading" class="loading">
        <div class="spinner"></div>
        <p>Loading mission data...</p>
      </div>

      <div *ngIf="error" class="error-message">
        {{ error }}
        <button (click)="error = null" class="close-error">✕</button>
      </div>

      <!-- No Mission State -->
      <div *ngIf="!mission && !loading" class="no-mission">
        <div class="no-mission-card">
          <h3>📋 No Active Mission</h3>
          <p>You currently have no assigned mission. Please wait for your next assignment from headquarters.</p>
          <div class="status-indicator">
            <span class="status-badge standby">🟡 STANDBY</span>
          </div>
        </div>
      </div>

      <!-- Mission Completed State -->
      <div *ngIf="mission && mission.is_completed && !loading" class="mission-completed">
        <div class="completed-card">
          <h3>🎉 Mission Completed!</h3>
          <p><strong>Mission:</strong> {{ mission.name }}</p>
          <p *ngIf="mission.completed_at"><strong>Completed on:</strong> {{ formatDate(mission.completed_at) }}</p>
          <p class="grateful-message">
            🌟 Excellent work, Agent {{ selectedCat.name }}! Your dedication and skill have successfully completed this mission. 
            All targets have been neutralized and objectives achieved. The agency is proud of your service.
          </p>
          <div class="status-indicator">
            <span class="status-badge completed">✅ COMPLETED</span>
          </div>
          <button class="btn btn-primary" (click)="returnToStandby()">
            🏠 Return to Standby
          </button>
        </div>
      </div>

      <!-- Active Mission State -->
      <div *ngIf="mission && !mission.is_completed && !loading" class="active-mission">
        <div class="mission-card">
          <div class="mission-header">
            <h3>🎯 {{ mission.name }}</h3>
            <div class="mission-status active">🟢 ACTIVE</div>
          </div>
          
          <div class="mission-content">
            <p class="mission-description">{{ mission.description }}</p>
            
            <div class="mission-dates" *ngIf="mission.start_date || mission.end_date">
              <div class="date-info" *ngIf="mission.start_date">
                <span class="date-label">Started:</span>
                <span class="date-value">{{ formatDate(mission.start_date) }}</span>
              </div>
              <div class="date-info" *ngIf="mission.end_date">
                <span class="date-label">Deadline:</span>
                <span class="date-value">{{ formatDate(mission.end_date) }}</span>
              </div>
            </div>

            <!-- Targets Management -->
            <div class="targets-section" *ngIf="mission.targets && mission.targets.length > 0">
              <h4>🎯 Mission Targets</h4>
              
              <div *ngFor="let target of mission.targets; let i = index" class="target-card">
                <div class="target-header">
                  <h5>{{ target.name }}</h5>
                  <span class="target-country">📍 {{ target.country }}</span>
                </div>
                
                <!-- Target Status -->
                <div class="target-status-section">
                  <label class="field-label">Status:</label>
                  <select 
                    [(ngModel)]="target.status" 
                    (change)="updateTargetStatus(target)"
                    [disabled]="isTargetFinal(target) || updating[target.id!]"
                    class="status-select"
                    [class]="'status-' + target.status">
                    <option value="init">📋 Init</option>
                    <option value="in_progress">🔄 In Progress</option>
                    <option value="completed">✅ Completed</option>
                  </select>
                  <div class="final-notice" *ngIf="isTargetFinal(target)">
                    🔒 Completed status is final and cannot be changed
                  </div>
                </div>

                <!-- Target Notes -->
                <div class="target-notes-section">
                  <label class="field-label">Field Notes:</label>
                  <textarea 
                    [(ngModel)]="target.notes"
                    (blur)="updateTargetNotes(target)"
                    [disabled]="isTargetFinal(target) || updating[target.id!]"
                    placeholder="Enter field observations, intelligence, and mission notes..."
                    class="notes-textarea"
                    rows="3">
                  </textarea>
                  <div class="final-notice" *ngIf="isTargetFinal(target)">
                    🔒 Target is completed and cannot be modified
                  </div>
                </div>

                <div class="target-meta">
                  <span class="meta-info">🕐 Updated: {{ formatDate(target.updated_at) }}</span>
                  <div *ngIf="updating[target.id!]" class="updating-indicator">
                    <span class="spinner-sm"></span> Updating...
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styleUrls: ['./cat-mission.component.css']
})
export class CatMissionComponent implements OnInit, OnChanges {
  @Input() selectedCat: SpyCat | null = null;
  
  mission: Mission | null = null;
  loading = false;
  error: string | null = null;
  updating: { [targetId: number]: boolean } = {};

  ngOnInit(): void {
    this.loadCatMission();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['selectedCat'] && this.selectedCat) {
      this.loadCatMission();
    }
  }

  loadCatMission(): void {
    if (!this.selectedCat) return;

    this.loading = true;
    this.error = null;
    
    fetch(`http://localhost:3001/api/v1/spy-cats/${this.selectedCat.id}/mission`)
      .then(response => {
        if (response.status === 204) {
          // No mission assigned
          this.mission = null;
          this.loading = false;
          return null;
        }
        if (!response.ok) {
          throw new Error('Failed to fetch mission');
        }
        return response.json();
      })
      .then(mission => {
        this.mission = mission;
        this.loading = false;
      })
      .catch(error => {
        console.error('Error loading cat mission:', error);
        this.error = 'Failed to load mission data. Please try again.';
        this.loading = false;
      });
  }

  updateTargetStatus(target: Target): void {
    if (!this.selectedCat || !target.id) return;

    this.updating[target.id] = true;
    this.error = null;

    fetch(`http://localhost:3001/api/v1/spy-cats/${this.selectedCat.id}/mission/targets/${target.id}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status: target.status })
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(err => Promise.reject(err));
        }
        return response.json();
      })
      .then(updatedTarget => {
        // Update target in the mission
        if (this.mission && this.mission.targets) {
          const index = this.mission.targets.findIndex(t => t.id === target.id);
          if (index !== -1) {
            this.mission.targets[index] = { ...this.mission.targets[index], ...updatedTarget };
          }
        }
        
        // Check if mission is now completed
        this.checkMissionCompletion();
        
        this.updating[target.id!] = false;
      })
      .catch(error => {
        console.error('Error updating target status:', error);
        this.error = error.details || 'Failed to update target status. Please try again.';
        this.updating[target.id!] = false;
        
        // Reload mission to get fresh data
        this.loadCatMission();
      });
  }

  updateTargetNotes(target: Target): void {
    if (!this.selectedCat || !target.id) return;

    this.updating[target.id] = true;
    this.error = null;

    fetch(`http://localhost:3001/api/v1/spy-cats/${this.selectedCat.id}/mission/targets/${target.id}/notes`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ notes: target.notes || '' })
    })
      .then(response => {
        if (!response.ok) {
          return response.json().then(err => Promise.reject(err));
        }
        return response.json();
      })
      .then(updatedTarget => {
        // Update target in the mission
        if (this.mission && this.mission.targets) {
          const index = this.mission.targets.findIndex(t => t.id === target.id);
          if (index !== -1) {
            this.mission.targets[index] = { ...this.mission.targets[index], ...updatedTarget };
          }
        }
        this.updating[target.id!] = false;
      })
      .catch(error => {
        console.error('Error updating target notes:', error);
        this.error = error.details || 'Failed to update target notes. Please try again.';
        this.updating[target.id!] = false;
        
        // Reload mission to get fresh data
        this.loadCatMission();
      });
  }

  isTargetFinal(target: Target): boolean {
    return target.status === 'completed';
  }

  checkMissionCompletion(): void {
    if (!this.mission || !this.mission.targets) return;

    const allCompleted = this.mission.targets.every(target => target.status === 'completed');
    if (allCompleted) {
      // Reload mission to get updated completion status
      setTimeout(() => {
        this.loadCatMission();
      }, 1000);
    }
  }

  returnToStandby(): void {
    // Reset mission and show standby state
    this.mission = null;
    // In a real app, you might want to emit an event to parent component
  }

  formatDate(date: Date | string | null | undefined): string {
    if (!date) return 'Not set';
    return new Date(date).toLocaleDateString();
  }
}
