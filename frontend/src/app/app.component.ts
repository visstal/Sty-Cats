import { Component, OnInit } from '@angular/core';

interface SpyCat {
  id: number;
  name: string;
  breed: string;
  yearsOfExperience: number;
  salary: number;
}

@Component({
  selector: 'app-root',
  template: `
    <div class="app-container">
      <app-floating-emojis></app-floating-emojis>
      
      <header class="app-header">
        <div class="header-content">
          <h1 class="logo">🕵️ SPY CAT AGENCY 🐱</h1>
          <nav class="nav-tabs">
            <button 
              class="nav-tab" 
              [class.active]="activeTab === 'agency'"
              (click)="setActiveTab('agency')">
              🏢 Agency
            </button>
            <button 
              class="nav-tab" 
              [class.active]="activeTab === 'spy-cats'"
              (click)="setActiveTab('spy-cats')">
              🐱 Spy Cats
            </button>
          </nav>
          
          <!-- Cat Selection Dropdown for Spy Cats mode -->
          <div class="cat-selector" *ngIf="activeTab === 'spy-cats'">
            <select [(ngModel)]="selectedCatId" (change)="onCatSelected()" class="cat-dropdown">
              <option value="">Select a Spy Cat...</option>
              <option *ngFor="let cat of cats" [value]="cat.id">
                🐱 {{ cat.name }} ({{ cat.breed }})
              </option>
            </select>
            <div class="welcome-message" *ngIf="selectedCat">
              Welcome, Agent {{ selectedCat.name }}! 🐾
            </div>
          </div>
        </div>
      </header>

      <main class="main-container">
        <!-- Agency Mode with Sub-tabs -->
        <div *ngIf="activeTab === 'agency'" class="agency-container">
          <div class="sub-nav">
            <button 
              class="sub-tab" 
              [class.active]="agencySubTab === 'cats'"
              (click)="setAgencySubTab('cats')">
              🐱 Cats Management
            </button>
            <button 
              class="sub-tab" 
              [class.active]="agencySubTab === 'missions'"
              (click)="setAgencySubTab('missions')">
              🎯 Missions Management
            </button>
          </div>
          
          <div class="sub-content">
            <app-cats *ngIf="agencySubTab === 'cats'"></app-cats>
            <app-missions *ngIf="agencySubTab === 'missions'"></app-missions>
          </div>
        </div>
        
        <div *ngIf="activeTab === 'spy-cats'" class="spy-mode-container">
          <div *ngIf="!selectedCat" class="select-cat-prompt">
            <h2>🕵️ Spy Cat Mode</h2>
            <p>Select your agent identity from the dropdown above to access your missions.</p>
          </div>
          
          <div *ngIf="selectedCat" class="cat-dashboard">
            <app-cat-mission [selectedCat]="selectedCat"></app-cat-mission>
          </div>
        </div>
      </main>
    </div>
  `,
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  activeTab = 'agency';
  agencySubTab = 'cats';
  cats: SpyCat[] = [];
  selectedCatId: string = '';
  selectedCat: SpyCat | null = null;

  ngOnInit() {
    this.loadCats();
  }

  setActiveTab(tab: string) {
    this.activeTab = tab;
    if (tab === 'agency') {
      this.agencySubTab = 'cats'; // Reset to cats when switching to agency
    }
    if (tab !== 'spy-cats') {
      this.selectedCatId = '';
      this.selectedCat = null;
    }
  }

  setAgencySubTab(subTab: string) {
    this.agencySubTab = subTab;
  }

  async loadCats() {
    try {
      const response = await fetch('http://localhost:3001/api/v1/cats');
      if (response.ok) {
        const data = await response.json();
        // Handle both direct array response and {cats: [...]} wrapper
        this.cats = Array.isArray(data) ? data : data.cats || [];
      }
    } catch (error) {
      console.error('Failed to load cats:', error);
    }
  }

  onCatSelected() {
    const foundCat = this.cats.find(cat => cat.id.toString() === this.selectedCatId);
    this.selectedCat = foundCat || null;
  }
}
