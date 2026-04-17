import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from "@/components/ui/sidebar";
import type { Database } from "@/types/types";
import {
  CaretDownIcon,
  DatabaseIcon,
  GearIcon,
  ShieldCheckIcon,
  type Icon,
  PlusSquareIcon,
} from "@phosphor-icons/react";
import { Link } from "@tanstack/react-router";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { DatabaseCreationDialog } from "@/components/database/database-creation-dialog.tsx";
import { useSidebar } from "@/components/ui/sidebar-context.tsx";

export type Submenu = {
  href: string;
  label: string;
  active?: boolean;
  icon?: Icon;
};

export type Menu = {
  href: string;
  label: string;
  active?: boolean;
  icon?: Icon;
  submenus?: Submenu[];
};

export type Group = {
  groupLabel?: string;
  menus: Menu[];
};

const useMenuItems = (projectId: string, databases?: Database[]): Group[] => {
  return [
    {
      menus: [
        {
          href: `/projects/${ projectId }/databases`,
          label: "Databases",
          icon: DatabaseIcon,
          submenus: databases?.map((db) => ({
            href: `/projects/${ projectId }/databases/${ db.id }`,
            label: db.displayName,
            icon: DatabaseIcon,
          })) as Submenu[],
        },
        {
          href: `/projects/${ projectId }/ci-cd`,
          label: "CI/CD Integration",
          icon: GearIcon,
          submenus: [],
        },
        {
          href: `/projects/${ projectId }/api-keys`,
          label: "API Keys",
          icon: ShieldCheckIcon,
          submenus: [],
        },
        {
          href: `/projects/${ projectId }/settings`,
          label: "Settings",
          icon: GearIcon,
          submenus: [],
        },
      ],
    },
  ];
};

export function AppSidebar({
                             projectId,
                             databases,
                           }: {
  projectId: string;
  databases?: Database[];
}) {
  const {setOpen} = useSidebar();
  const groups = useMenuItems(projectId, databases);
  return (
    <Sidebar variant="inset" collapsible="icon">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>
            <Link to="/">Skemr</Link>
          </SidebarGroupLabel>
          <SidebarGroupContent>
            { groups.map(({groupLabel, menus}, index) => (
              <SidebarGroup key={ groupLabel || index }>
                { groupLabel && (
                  <SidebarGroupLabel>{ groupLabel }</SidebarGroupLabel>
                ) }
                <SidebarGroupContent>
                  <SidebarMenu>
                    { menus.map(
                      ({href, label, icon: Icon, active, submenus}) =>
                        submenus?.length === 0 ? (
                          <SidebarMenuItem key={ href }>
                            { label === "Databases" ? (
                              <Collapsible>
                                <CollapsibleTrigger
                                  onClick={ () => setOpen(true) }
                                  render={
                                    <SidebarMenuButton
                                      isActive={ active }
                                      tooltip={ label }
                                    >
                                      { Icon && <Icon size={ 18 }/> }
                                      <span>{ label }</span>
                                      <CaretDownIcon
                                        className={ `ml-auto transition-transform group-data-open/collapsible:rotate-180` }
                                      />
                                    </SidebarMenuButton>
                                  }
                                ></CollapsibleTrigger>
                                <CollapsibleContent>
                                  <SidebarMenuSub>
                                    <DatabaseCreationDialog
                                      projectId={ projectId }
                                      trigger={
                                        <SidebarMenuSubButton>
                                          <PlusSquareIcon size={ 16 }/> Add
                                          Database
                                        </SidebarMenuSubButton>
                                      }
                                    />
                                  </SidebarMenuSub>
                                </CollapsibleContent>
                              </Collapsible>
                            ) : (
                              <SidebarMenuButton
                                tooltip={ label }
                                render={
                                  <Link to={ href }>
                                    { Icon && <Icon size={ 18 }/> }
                                    <span>{ label }</span>
                                  </Link>
                                }
                                isActive={ active }
                              />
                            ) }
                          </SidebarMenuItem>
                        ) : (
                          <SidebarMenuItem key={ `group-${ label }` }>
                            <Collapsible className="group/collapsible">
                              <CollapsibleTrigger
                                onClick={ () => setOpen(true) }
                                render={
                                  <SidebarMenuButton
                                    isActive={ active }
                                    tooltip={ label }
                                  >
                                    { Icon && <Icon size={ 18 }/> }
                                    <span>{ label }</span>
                                    <CaretDownIcon
                                      className={ `ml-auto transition-transform group-data-open/collapsible:rotate-180` }
                                    />
                                  </SidebarMenuButton>
                                }
                              ></CollapsibleTrigger>

                              <CollapsibleContent>
                                <SidebarMenuSub>
                                  { submenus?.map(
                                    ({
                                       href,
                                       label,
                                       icon: SubIcon,
                                       active,
                                     }) => (
                                      <SidebarMenuSubItem key={ href }>
                                        <SidebarMenuSubButton
                                          render={
                                            <Link to={ href }>
                                              { SubIcon && <SubIcon size={ 18 }/> }
                                              <span>{ label }</span>
                                            </Link>
                                          }
                                          isActive={ active }
                                        />
                                      </SidebarMenuSubItem>
                                    ),
                                  ) }
                                </SidebarMenuSub>

                                { label === "Databases" && (
                                  <SidebarMenuSub>
                                    <DatabaseCreationDialog
                                      projectId={ projectId }
                                      trigger={
                                        <SidebarMenuSubButton>
                                          <PlusSquareIcon size={ 16 }/> Add
                                          Database
                                        </SidebarMenuSubButton>
                                      }
                                    />
                                  </SidebarMenuSub>
                                ) }
                              </CollapsibleContent>
                            </Collapsible>
                          </SidebarMenuItem>
                        ),
                    ) }
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            )) }
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
