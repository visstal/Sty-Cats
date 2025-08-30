import { Component, OnInit } from '@angular/core';
import { Mission, CreateMissionRequest, CreateTargetRequest } from '../interfaces/mission.interface';
import { MissionService } from '../services/mission.service';

@Component({
  selector: 'app-missions',
  templateUrl: './missions.component.html',
  styleUrls: ['./missions.component.css']
})
export class MissionsComponent implements OnInit {
  missions: Mission[] = [];
  loading = false;
  error: string | null = null;
  showCreateForm = false;
  showAssignDialog = false;
  selectedMission: Mission | null = null;
  selectedCatId: string | number | null = null;
  freeCats: any[] = [];
  
  // Target management for existing missions
  showTargetManagement: { [missionId: number]: boolean } = {};
  newTargetData: { [missionId: number]: { name: string; country: string } } = {};
  
  // Notes popup management
  showNotes: { [targetId: number]: boolean } = {};
  
  newMission: CreateMissionRequest = {
    name: '',
    description: '',
    targets: []
  };

  // Target management
  targetName = '';
  targetCountry = '';

  constructor(private missionService: MissionService) { }

  ngOnInit(): void {
    this.loadMissions();
  }

  loadMissions(): void {
    this.loading = true;
    this.error = null;
    
    this.missionService.getMissions().subscribe({
      next: (missions: Mission[]) => {
        this.missions = missions;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading missions:', error);
        this.error = 'Failed to load missions. Please try again.';
        this.loading = false;
      }
    });
  }

  toggleCreateForm(): void {
    this.showCreateForm = !this.showCreateForm;
    if (!this.showCreateForm) {
      this.resetForm();
    }
  }

  createMission(): void {
    if (this.isFormValid()) {
      this.loading = true;
      
      // Create mission payload including targets (required)
      const missionData: CreateMissionRequest = {
        name: this.newMission.name,
        description: this.newMission.description,
        targets: this.newMission.targets // Now required, validated by isFormValid()
      };
      
      this.missionService.createMission(missionData).subscribe({
        next: (mission) => {
          this.missions.push(mission);
          this.resetForm();
          this.showCreateForm = false;
          this.loading = false;
        },
        error: (error) => {
          console.error('Error creating mission:', error);
          if (error.error && error.error.details) {
            this.error = `Failed to create mission: ${error.error.details}`;
          } else {
            this.error = 'Failed to create mission. Please try again.';
          }
          this.loading = false;
        }
      });
    }
  }

  deleteMission(mission: Mission): void {
    if (mission.cat) {
      alert('Cannot delete mission with assigned agent. Please unassign the agent first.');
      return;
    }
    
    if (mission.id && confirm(`Are you sure you want to delete "${mission.name}"?`)) {
      this.loading = true;
      
      this.missionService.deleteMission(mission.id).subscribe({
        next: () => {
          this.missions = this.missions.filter(m => m.id !== mission.id);
          this.loading = false;
        },
        error: (error) => {
          console.error('Error deleting mission:', error);
          this.error = 'Failed to delete mission. Please try again.';
          this.loading = false;
        }
      });
    }
  }

  public isFormValid(): boolean {
    return !!(
      this.newMission.name?.trim() &&
      this.newMission.description?.trim() &&
      this.newMission.targets &&
      this.newMission.targets.length >= 1 &&
      this.newMission.targets.length <= 3
    );
  }

  private resetForm(): void {
    this.newMission = {
      name: '',
      description: '',
      start_date: '',
      end_date: '',
      targets: []
    };
    this.targetName = '';
    this.targetCountry = '';
    this.error = null;
  }

  // Target management methods
  addTarget(): void {
    if (this.targetName.trim() && this.targetCountry.trim() && (this.newMission.targets?.length || 0) < 3) {
      if (!this.newMission.targets) {
        this.newMission.targets = [];
      }
      this.newMission.targets.push({
        name: this.targetName.trim(),
        country: this.targetCountry.trim()
      });
      
      // Clear only target input fields
      this.targetName = '';
      this.targetCountry = '';
    }
  }

  removeTarget(index: number): void {
    if (this.newMission.targets) {
      this.newMission.targets.splice(index, 1);
    }
  }

  canAddTarget(): boolean {
    return this.targetName.trim() !== '' && 
           this.targetCountry.trim() !== '' && 
           (this.newMission.targets?.length || 0) < 3;
  }

  formatDate(date: Date | string | null | undefined): string {
    if (!date) {
      return 'Not set';
    }
    return new Date(date).toLocaleDateString();
  }

  getMissionStatus(mission: Mission): string {
    if (mission.is_completed) {
      return 'Completed';
    }
    
    if (!mission.end_date) {
      return 'In Progress';
    }
    
    const endDate = new Date(mission.end_date);
    const now = new Date();
    
    if (endDate < now) {
      return 'Overdue';
    }
    
    if (!mission.start_date) {
      return 'Pending';
    }
    
    const startDate = new Date(mission.start_date);
    if (startDate <= now && endDate >= now) {
      return 'Active';
    }
    
    return 'Pending';
  }

  openAssignDialog(mission: Mission): void {
    this.selectedMission = mission;
    this.selectedCatId = null;
    this.showAssignDialog = true;
    this.loadFreeCats();
  }

  closeAssignDialog(): void {
    this.showAssignDialog = false;
    this.selectedMission = null;
    this.selectedCatId = null;
    this.freeCats = [];
  }

  loadFreeCats(): void {
    this.missionService.getFreeCats().subscribe({
      next: (cats) => {
        this.freeCats = cats;
      },
      error: (error) => {
        console.error('Error loading free cats:', error);
        this.error = 'Failed to load available cats';
      }
    });
  }

  assignCat(): void {
    if (!this.selectedMission || !this.selectedCatId) {
      return;
    }

    this.loading = true;
    this.missionService.assignCatToMission(this.selectedMission.id!, +this.selectedCatId).subscribe({
      next: (updatedMission) => {
        // Update the mission in the list
        const index = this.missions.findIndex(m => m.id === updatedMission.id);
        if (index !== -1) {
          this.missions[index] = updatedMission;
        }
        this.closeAssignDialog();
        this.loading = false;
      },
      error: (error) => {
        console.error('Error assigning cat:', error);
        let errorMessage = 'Failed to assign cat to mission. Please try again.';
        
        // Extract detailed validation messages from backend response
        if (error.status === 400 && error.error) {
          if (error.error.details) {
            errorMessage = `${error.error.details}`;
          } else if (error.error.error) {
            errorMessage = `${error.error.error}`;
          }
        } else if (error.error && error.error.details) {
          errorMessage = `${error.error.details}`;
        }
        
        this.error = errorMessage;
        this.loading = false;
      }
    });
  }

  getStatusClass(mission: Mission): string {
    const status = this.getMissionStatus(mission);
    return status.toLowerCase();
  }

  // Target management methods for existing missions
  toggleTargetManagement(missionId: number): void {
    this.showTargetManagement[missionId] = !this.showTargetManagement[missionId];
    if (!this.newTargetData[missionId]) {
      this.newTargetData[missionId] = { name: '', country: '' };
    }
  }

  canAddTargetToMission(mission: Mission): boolean {
    return (mission.targets?.length || 0) < 3;
  }

  canDeleteTarget(target: any): boolean {
    return target.status === 'init';
  }

  isValidTargetData(missionId: number): boolean {
    const targetData = this.newTargetData[missionId];
    return targetData && targetData.name.trim() !== '' && targetData.country.trim() !== '';
  }

  addTargetToExistingMission(missionId: number): void {
    if (!this.isValidTargetData(missionId)) {
      return;
    }

    const targetData = this.newTargetData[missionId];
    this.loading = true;

    this.missionService.addTargetToMission(missionId, {
      name: targetData.name.trim(),
      country: targetData.country.trim()
    }).subscribe({
      next: (newTarget) => {
        // Find the mission and add the new target
        const mission = this.missions.find(m => m.id === missionId);
        if (mission) {
          if (!mission.targets) {
            mission.targets = [];
          }
          mission.targets.push(newTarget);
        }
        
        // Reset form
        this.newTargetData[missionId] = { name: '', country: '' };
        this.loading = false;
      },
      error: (error) => {
        console.error('Error adding target:', error);
        if (error.error && error.error.details) {
          this.error = `Failed to add target: ${error.error.details}`;
        } else {
          this.error = 'Failed to add target. Please try again.';
        }
        this.loading = false;
      }
    });
  }

  deleteTargetFromMission(missionId: number, targetId: number, target: any): void {
    if (!this.canDeleteTarget(target)) {
      alert('Only targets in "init" status can be deleted.');
      return;
    }

    if (confirm(`Are you sure you want to delete target "${target.name}"?`)) {
      this.loading = true;

      this.missionService.deleteTargetFromMission(missionId, targetId).subscribe({
        next: () => {
          // Remove the target from the mission
          const mission = this.missions.find(m => m.id === missionId);
          if (mission && mission.targets) {
            mission.targets = mission.targets.filter(t => t.id !== targetId);
          }
          this.loading = false;
        },
        error: (error) => {
          console.error('Error deleting target:', error);
          if (error.error && error.error.details) {
            this.error = `Failed to delete target: ${error.error.details}`;
          } else {
            this.error = 'Failed to delete target. Please try again.';
          }
          this.loading = false;
        }
      });
    }
  }

  toggleNotes(targetId: number): void {
    this.showNotes[targetId] = !this.showNotes[targetId];
  }
}
