import AdminLayout from './layouts/AdminLayout.svelte';
import ErrorLayout from './layouts/ErrorPage.svelte';
import IndexLayout from './layouts/IndexLayout.svelte';
import ManageLayout from './layouts/ManageLayout.svelte';
import TranscriptViewLayout from './layouts/TranscriptViewLayout.svelte';

import BotStaff from './views/admin/BotStaff.svelte';
import Appearance from './views/Appearance.svelte';
import Blacklist from './views/Blacklist.svelte';
import Error from './views/Error.svelte';
import Error404 from './views/Error404.svelte';
import Forms from './views/Forms.svelte';
import Import from './views/Import.svelte';
import Index from './views/Index.svelte';
import IntegrationActivate from './views/integrations/Activate.svelte';
import IntegrationConfigure from './views/integrations/Configure.svelte';
import IntegrationCreate from './views/integrations/Create.svelte';
import Integrations from './views/integrations/Integrations.svelte';
import IntegrationManage from './views/integrations/Manage.svelte';
import IntegrationView from './views/integrations/View.svelte';
import Login from './views/Login.svelte';
import LoginCallback from './views/LoginCallback.svelte';
import Logout from './views/Logout.svelte';
import CreateMultiPanel from './views/panels/CreateMultiPanel.svelte';
import CreatePanel from './views/panels/CreatePanel.svelte';
import EditMultiPanel from './views/panels/EditMultiPanel.svelte';
import EditPanel from './views/panels/EditPanel.svelte';
import Panels from './views/panels/Panels.svelte';
import SelectServers from './views/premium/SelectServers.svelte';
import Settings from './views/Settings.svelte';
import StaffOverride from './views/StaffOverride.svelte';
import Tags from './views/Tags.svelte';
import Teams from './views/Teams.svelte';
import Tickets from './views/Tickets.svelte';
import TicketView from './views/TicketView.svelte';
import Transcripts from './views/Transcripts.svelte';
import TranscriptView from './views/TranscriptView.svelte';
import Whitelabel from './views/Whitelabel.svelte';

export const routes = [
	{ name: '/', component: Index, layout: IndexLayout },
	{ name: '404', path: '404', component: Error404, layout: ErrorLayout },
	{ name: '/callback', component: LoginCallback },
	{ name: '/login', component: Login },
	{ name: '/logout', component: Logout },
	{ name: '/error', component: Error, layout: ErrorLayout },
	{ name: '/whitelabel', component: Whitelabel, layout: IndexLayout },
	{
		name: 'premium',
		nestedRoutes: [{ name: 'select-servers', component: SelectServers, layout: IndexLayout }],
	},
	{
		name: 'admin',
		nestedRoutes: [{ name: 'bot-staff', component: BotStaff, layout: AdminLayout }],
	},
	{
		name: 'manage/:id',
		nestedRoutes: [
			{ name: 'index', component: Error404, layout: ErrorLayout },
			{ name: 'settings', component: Settings, layout: ManageLayout },
			{ name: 'appearance', component: Appearance, layout: ManageLayout },
			{
				name: 'transcripts',
				nestedRoutes: [
					{
						name: 'index',
						component: Transcripts,
						layout: ManageLayout,
					},
					{
						name: 'view/:ticketid',
						component: TranscriptView, // just to test
						layout: TranscriptViewLayout,
					},
				],
			},
			// Backwards compatibility
			{
				name: 'logs',
				nestedRoutes: [
					{
						name: 'view/:ticketid',
						component: TranscriptView,
						layout: TranscriptViewLayout,
					},
				],
			},
			{
				name: 'panels',
				nestedRoutes: [
					{
						name: 'index',
						component: Panels,
						layout: ManageLayout,
					},
					{
						name: 'create',
						component: CreatePanel,
						layout: ManageLayout,
					},
					{
						name: 'create-multi',
						component: CreateMultiPanel,
						layout: ManageLayout,
					},
					{
						name: 'edit/:panelid',
						component: EditPanel,
						layout: ManageLayout,
					},
					{
						name: 'edit-multi/:panelid',
						component: EditMultiPanel,
						layout: ManageLayout,
					},
				],
			},
			{ name: 'blacklist', component: Blacklist, layout: ManageLayout },
			{ name: 'tags', component: Tags, layout: ManageLayout },
			{ name: 'import', component: Import, layout: ManageLayout },
			{ name: 'teams', component: Teams, layout: ManageLayout },
			{ name: 'forms', component: Forms, layout: ManageLayout },
			{ name: 'staffoverride', component: StaffOverride, layout: ManageLayout },
			{
				name: 'tickets',
				nestedRoutes: [
					{
						name: 'index',
						component: Tickets,
						layout: ManageLayout,
					},
					{
						name: 'view/:ticketid',
						component: TicketView,
						layout: ManageLayout,
					},
				],
			},
			{
				name: 'integrations',
				nestedRoutes: [
					{
						name: 'index',
						component: Integrations,
						layout: ManageLayout,
					},
					{
						name: 'create',
						component: IntegrationCreate,
						layout: ManageLayout,
					},
					{
						name: '/view/:integration',
						component: IntegrationView,
						layout: ManageLayout,
					},
					{
						name: '/configure/:integration',
						component: IntegrationConfigure,
						layout: ManageLayout,
					},
					{
						name: '/activate/:integration',
						component: IntegrationActivate,
						layout: ManageLayout,
					},
					{
						name: '/manage/:integration',
						component: IntegrationManage,
						layout: ManageLayout,
					},
				],
			},
		],
	},
];
