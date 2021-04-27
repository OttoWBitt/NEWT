import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LandingComponent } from './views/landing/landing.component';
import { LandingRoutingModule } from './views/landing/landing-routing.module';
import { SystemLayoutComponent } from './shared/system-layout/system-layout.component';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { AuthInterceptor } from './services/auth/auth.interceptor';
import { ArtifactsComponent } from './shared/artifacts/artifacts.component';

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
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
