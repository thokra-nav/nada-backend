import styled from "styled-components";
import Link from 'next/link'
import {navLysGra} from "../../styles/constants";
import {components} from "../../lib/schema"

const SearchResultDiv = styled.div`
    background-color: ${navLysGra};
    padding: 5px;
    margin-bottom: 5px;
    h1 {
        font-size: 1.5em;
        margin: 0;
    }
`


type SearchResultEntry = components["schemas"]["SearchResultEntry"]

export interface SearchResultProps {
    searchResultEntry: SearchResultEntry
}

export const SearchResult = ({searchResultEntry}: SearchResultProps) => (
    <SearchResultDiv>
        <Link href={`/dataproduct/${searchResultEntry.id}`}>
            <h1>{searchResultEntry.name}</h1>
        </Link>
        <p>{searchResultEntry.excerpt}</p>
    </SearchResultDiv>
)

export default SearchResult