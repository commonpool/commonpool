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
import {ResourceInquiryComponent} from './resources/resource-inquiry/resource-inquiry.component';
import {OfferListComponent} from './offers/offer-list/offer-list.component';
import {CreateOfferComponent} from './offers/create-offer/create-offer.component';
import {OfferDetailsComponent} from './offers/offer-details/offer-details.component';
import {UserPickerComponent} from './shared/user-picker/user-picker.component';
import {ResourcePickerComponent} from './shared/resource-picker/resource-picker.component';
import {NgSelectModule} from '@ng-select/ng-select';
import {PlusSquareIcon} from './icons/plus-square/plus-square.icon';
import {PlusIcon} from './icons/plus/plus.icon';
import {UsernameComponent} from './shared/username/username.component';
import {ResourceNameComponent} from './shared/resource-name/resource-name.component';
import {TrashIcon} from './icons/trash/trash.icon';
import {CreateOrEditGroupComponent} from './groups/create-or-edit-group/create-or-edit-group.component';
import {RequiredIndicatorComponent} from './shared/required-indicator/required-indicator.component';
import {GroupViewComponent} from './groups/group-view/group-view.component';
import {BoxSeamIcon} from './icons/box-seam/box-seam.icon';
import {PersonIcon} from './icons/person/person.icon';
import {PeopleIcon} from './icons/people/people.icon';
import {PentagonIcon} from './icons/pentagon/pentagon.icon';
import {AsterixIcon} from './icons/asterix/asterix.icon';
import {GroupResourcesViewComponent} from './groups/group-resources-view/group-resources-view.component';
import {GroupMembersViewComponent} from './groups/group-members-view/group-members-view.component';

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
    OfferListComponent,
    CreateOfferComponent,
    OfferDetailsComponent,
    UserPickerComponent,
    ResourcePickerComponent,
    PlusSquareIcon,
    PlusIcon,
    UsernameComponent,
    ResourceNameComponent,
    TrashIcon,
    CreateOrEditGroupComponent,
    RequiredIndicatorComponent,
    GroupViewComponent,
    BoxSeamIcon,
    PersonIcon,
    PeopleIcon,
    PentagonIcon,
    AsterixIcon,
    GroupResourcesViewComponent,
    GroupMembersViewComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    ScrollingModule,
    NgSelectModule
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: AppHttpInterceptor, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(authService: AuthService) {
    authService.checkLoggedIn();
  }
}
