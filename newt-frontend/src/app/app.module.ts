import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { ArtifactsComponent } from './shared/artifacts/artifacts.component';
import { SystemLayoutComponent } from './shared/system-layout/system-layout.component';
import { LandingRoutingModule } from './views/landing/landing-routing.module';
import { LandingComponent } from './views/landing/landing.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatButtonModule } from '@angular/material/button';
import {MatDividerModule} from '@angular/material/divider';

@NgModule({
  declarations: [
    AppComponent,
    LandingComponent,
    SystemLayoutComponent,
    ArtifactsComponent
  ],
  imports: [
    BrowserModule,
    LandingRoutingModule,
    AppRoutingModule,
    FlexLayoutModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatDividerModule,
    BrowserAnimationsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
