import {
    expandAllFeature,
    hotkeysCoreFeature,
    searchFeature,
    selectionFeature,
    syncDataLoaderFeature,
    type TreeState,
} from "@headless-tree/core";
import { useTree } from "@headless-tree/react";
import {
    ArticleIcon,
    CaretUpDownIcon,
    CheckIcon,
    CircleNotchIcon,
    FunnelIcon,
    GraphIcon,
    MagnifyingGlassIcon,
    ShieldCheckIcon,
    TableIcon,
    XCircleIcon,
} from "@phosphor-icons/react";
import type React from "react";
import { useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import {
    InputGroup,
    InputGroupAddon,
    InputGroupButton,
    InputGroupInput,
} from "@/components/ui/input-group";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { Tree, TreeItem, TreeItemLabel } from "@/components/ui/tree";
import type { DatabaseEntity, DatabaseEntityType } from "@/types/types";
import { cn } from "@/lib/utils";

interface EntityItem {
    name: string;
    type: DatabaseEntityType | "root";
    children?: string[];
}

const ROOT_ID = "__root__";
const INDENT = 16;

function buildEntityRecord(
    entities: DatabaseEntity[],
): Record<string, EntityItem> {
    const record: Record<string, EntityItem> = {
        [ROOT_ID]: { name: "Root", type: "root", children: [] },
    };

    const schemas = entities.filter((e) => e.type === "schema");
    const tables = entities.filter((e) => e.type === "table");
    const columns = entities.filter((e) => e.type === "column");
    const constraints = entities.filter((e) => e.type === "constraint");
    const indexes = entities.filter((e) => e.type === "index");

    record[ROOT_ID]!.children = schemas.map((s) => s.id);

    for (const schema of schemas) {
        const schemaTableIds = tables
            .filter((t) => t.parentId === schema.id)
            .map((t) => t.id);
        record[schema.id] = {
            name: schema.name,
            type: "schema",
            children: schemaTableIds,
        };
    }

    for (const table of tables) {
        const tableColumnIds = () => {
            const tableColumns = columns
                .filter((c) => c.parentId === table.id)
                .map((c) => c.id);
            const tableIndexes = indexes
                .filter((i) => i.parentId === table.id)
                .map((i) => i.id);
            const tableConstraints = constraints
                .filter((c) => c.parentId === table.id)
                .map((c) => c.id);
            return [...tableColumns, ...tableIndexes, ...tableConstraints];
        };
        record[table.id] = {
            name: table.name,
            type: "table",
            children: tableColumnIds(),
        };
    }

    for (const column of columns) {
        record[column.id] = { name: column.name, type: "column" };
    }

    for (const constraint of constraints) {
        record[constraint.id] = { name: constraint.name, type: "constraint" };
    }

    for (const index of indexes) {
        record[index.id] = { name: index.name, type: "index" };
    }

    return record;
}

interface EntityTreeSelectorProps {
    entities: DatabaseEntity[];
    value: string;
    onChange: (id: string) => void;
    onPopoverOpenChange?: (open: boolean) => void;
    popoverOpen: boolean;
    loading?: boolean;
    id?: string;
}

const getIconByEntityType = (type: DatabaseEntityType | "root") => {
    switch (type) {
        case "table":
            return <TableIcon />;
        case "column":
            return <ArticleIcon />;
        case "schema":
            return <GraphIcon />;
        case "constraint":
            return <ShieldCheckIcon />;
        case "index":
            return <MagnifyingGlassIcon />;
        default:
            return null;
    }
};

export function EntityTreeSelector({
    entities,
    value,
    onChange,
    onPopoverOpenChange,
    popoverOpen,
    loading,
    id,
}: EntityTreeSelectorProps) {
    const selectedEntity = entities.find((e) => e.id === value);

    const entityRecord = useMemo(() => buildEntityRecord(entities), [entities]);
    const initialExpandedItems = useMemo(
        () =>
            Object.keys(entityRecord).filter(
                (id) => id !== ROOT_ID && entityRecord[id]?.type === "schema",
            ),
        [entityRecord],
    );

    const [state, setState] = useState<Partial<TreeState<EntityItem>>>({});
    const [searchValue, setSearchValue] = useState("");
    const [filteredItems, setFilteredItems] = useState<string[]>([]);
    const inputRef = useRef<HTMLInputElement>(null);

    const tree = useTree<EntityItem>({
        dataLoader: {
            getChildren: (itemId) => entityRecord[itemId]?.children ?? [],
            getItem: (itemId) =>
                entityRecord[itemId] ?? { name: itemId, type: "column" },
        },
        features: [
            syncDataLoaderFeature,
            hotkeysCoreFeature,
            selectionFeature,
            searchFeature,
            expandAllFeature,
        ],
        getItemName: (item) => item.getItemData().name,
        indent: INDENT,
        initialState: {
            expandedItems: initialExpandedItems,
        },
        isItemFolder: (item) => (item.getItemData()?.children?.length ?? 0) > 0,
        rootItemId: ROOT_ID,
        setState,
        state,
    });

    const handleClearSearch = () => {
        setSearchValue("");
        const searchProps = tree.getSearchInputElementProps();
        if (searchProps.onChange) {
            const syntheticEvent = {
                target: { value: "" },
            } as React.ChangeEvent<HTMLInputElement>;
            searchProps.onChange(syntheticEvent);
        }
        setState((prevState) => ({
            ...prevState,
            expandedItems: initialExpandedItems,
        }));
        setFilteredItems([]);
        if (inputRef.current) {
            inputRef.current.focus();
            inputRef.current.value = "";
        }
    };

    const handleClick = (itemId: string) => {
        onChange(itemId);
        onPopoverOpenChange?.(false);
    };

    const shouldShowItem = (itemId: string) => {
        if (!searchValue || searchValue.length === 0) return true;
        return filteredItems.includes(itemId);
    };

    useEffect(() => {
        if (!searchValue || searchValue.length === 0) {
            setFilteredItems([]);
            return;
        }
        const allItems = tree.getItems();

        const directMatches = allItems
            .filter((item) =>
                item
                    .getItemName()
                    .toLowerCase()
                    .includes(searchValue.toLowerCase()),
            )
            .map((item) => item.getId());

        const parentIds = new Set<string>();
        for (const matchId of directMatches) {
            let item = allItems.find((i) => i.getId() === matchId);
            while (item?.getParent?.()) {
                const parent = item.getParent();
                if (parent) {
                    parentIds.add(parent.getId());
                    item = parent;
                } else {
                    break;
                }
            }
        }

        const childrenIds = new Set<string>();
        for (const matchId of directMatches) {
            const item = allItems.find((i) => i.getId() === matchId);
            if (item?.isFolder()) {
                const getDescendants = (itemId: string) => {
                    const children = entityRecord[itemId]?.children || [];
                    for (const childId of children) {
                        childrenIds.add(childId);
                        if (entityRecord[childId]?.children?.length) {
                            getDescendants(childId);
                        }
                    }
                };
                getDescendants(item.getId());
            }
        }

        setFilteredItems([
            ...directMatches,
            ...Array.from(parentIds),
            ...Array.from(childrenIds),
        ]);

        const folderIdsToExpand = allItems
            .filter((item) => item.isFolder())
            .map((item) => item.getId());

        setState((prevState) => ({
            ...prevState,
            expandedItems: [
                ...new Set([
                    ...(tree.getState().expandedItems || []),
                    ...folderIdsToExpand,
                ]),
            ],
        }));
    }, [searchValue, tree, entityRecord]);

    return (
        <Popover open={popoverOpen} onOpenChange={onPopoverOpenChange}>
            <PopoverTrigger
                render={
                    <Button
                        id={id}
                        variant="outline"
                        role="combobox"
                        aria-expanded={popoverOpen}
                        className="w-full justify-between font-normal"
                        disabled={loading}
                    />
                }
            >
                {loading ? (
                    <CircleNotchIcon className="size-4 animate-spin text-muted-foreground" />
                ) : (
                    <span className="truncate">
                        {selectedEntity ? selectedEntity.name : "Select entity"}
                    </span>
                )}
                <CaretUpDownIcon className="ml-2 shrink-0 opacity-50" />
            </PopoverTrigger>
            <PopoverContent
                className="w-72 flex flex-col gap-2 p-1 "
                align="start"
                side="bottom"
            >
                {/* Search input */}
                <InputGroup>
                    <InputGroupAddon align="inline-start">
                        <FunnelIcon aria-hidden="true" className="size-3.5" />
                    </InputGroupAddon>
                    <InputGroupInput
                        ref={inputRef}
                        className="text-sm"
                        onBlur={(e) => {
                            e.preventDefault();
                            if (searchValue && searchValue.length > 0) {
                                const searchProps =
                                    tree.getSearchInputElementProps();
                                if (searchProps.onChange) {
                                    const syntheticEvent = {
                                        target: { value: searchValue },
                                    } as React.ChangeEvent<HTMLInputElement>;
                                    searchProps.onChange(syntheticEvent);
                                }
                            }
                        }}
                        onChange={(e) => {
                            const val = e.target.value;
                            setSearchValue(val);
                            const searchProps =
                                tree.getSearchInputElementProps();
                            if (searchProps.onChange) {
                                searchProps.onChange(e);
                            }
                            if (val.length > 0) {
                                tree.expandAll();
                            } else {
                                setState((prevState) => ({
                                    ...prevState,
                                    expandedItems: initialExpandedItems,
                                }));
                                setFilteredItems([]);
                            }
                        }}
                        placeholder="Filter entities..."
                        type="search"
                        value={searchValue}
                    />
                    {searchValue && (
                        <InputGroupAddon align="inline-end">
                            <InputGroupButton
                                aria-label="Clear search"
                                onClick={handleClearSearch}
                                size="icon-xs"
                                type="button"
                                variant="ghost"
                            >
                                <XCircleIcon
                                    aria-hidden="true"
                                    className="size-3.5"
                                />
                            </InputGroupButton>
                        </InputGroupAddon>
                    )}
                </InputGroup>

                {/* Tree */}
                <div className="max-h-60 pr-2 overflow-y-auto">
                    <Tree indent={INDENT} tree={tree}>
                        {searchValue && filteredItems.length === 0 ? (
                            <p className="px-3 py-4 text-center text-sm text-muted-foreground">
                                No entities found for "{searchValue}"
                            </p>
                        ) : (
                            tree.getItems().map((item) => {
                                const isVisible = shouldShowItem(item.getId());
                                const entityType = item.getItemData().type;
                                const isSelected = item.getId() === value;

                                return (
                                    <TreeItem
                                        className="data-[visible=false]:hidden w-full text-left"
                                        data-visible={isVisible || !searchValue}
                                        item={item}
                                        key={item.getId()}
                                    >
                                        <TreeItemLabel
                                            className={cn(
                                                isSelected
                                                    ? "bg-accent text-accent-foreground"
                                                    : "",
                                                "border border-border px-0.5",
                                            )}
                                        >
                                            <span className="flex items-center  gap-2 w-full">
                                                <span className="size-3.5 shrink-0 text-muted-foreground">
                                                    {getIconByEntityType(
                                                        entityType,
                                                    )}
                                                </span>
                                                <span className="truncate">
                                                    {item.getItemName()}
                                                </span>
                                                <Button
                                                    onClick={() => {
                                                        handleClick(
                                                            item.getId(),
                                                        );
                                                    }}
                                                    className="ml-auto"
                                                    variant={
                                                        isSelected
                                                            ? "default"
                                                            : "outline"
                                                    }
                                                    size={"icon-sm"}
                                                >
                                                    <CheckIcon
                                                        className={cn(
                                                            isSelected
                                                                ? "opacity-100"
                                                                : "opacity-50",
                                                        )}
                                                    />
                                                </Button>
                                            </span>
                                        </TreeItemLabel>
                                    </TreeItem>
                                );
                            })
                        )}
                    </Tree>
                </div>
            </PopoverContent>
        </Popover>
    );
}
