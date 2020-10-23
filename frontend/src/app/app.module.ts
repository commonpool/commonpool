import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import {TopNavComponent} from './top-nav/top-nav.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {AuthService} from './auth.service';
import {AppHttpInterceptor} from './api/backend.service';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
import {UserProfileComponent} from './user-profile/user-profile.component';
import {ConversationThreadComponent} from './chat/conversation-thread/conversation-thread.component';
import {ConversationThreadListComponent} from './chat/conversation-thread-list/conversation-thread-list.component';
import {ScrollingModule} from '@angular/cdk/scrolling';
import { ResourceInquiryComponent } from './resources/resource-inquiry/resource-inquiry.component';

@NgModule({
  declarations: [
    AppComponent,
    ResourceListViewComponent,
    TopNavComponent,
    CreateOrEditResourceComponent,
    ResourceDetailsComponent,
    UserProfileComponent,
    ConversationThreadComponent,
    ConversationThreadListComponent,
    ResourceInquiryComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    ScrollingModule
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: AppHttpInterceptor, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(authService: AuthService) {
    authService.checkLoggedIn();
  }
}
