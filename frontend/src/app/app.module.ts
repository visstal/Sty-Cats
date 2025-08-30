import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { CatsComponent } from './cats/cats.component';
import { MissionsComponent } from './missions/missions.component';
import { FloatingEmojisComponent } from './floating-emojis/floating-emojis.component';
import { CatMissionComponent } from './cat-mission/cat-mission.component';

@NgModule({
  declarations: [
    AppComponent,
    CatsComponent,
    MissionsComponent,
    FloatingEmojisComponent,
    CatMissionComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
